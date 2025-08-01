# Proto to JSON Conversion Guide

This document explains the improved proto to JSON conversion system in protoc-gen-go-wasmjs, particularly focusing on handling differences between Go's protojson and TypeScript protobuf libraries.

## Overview

The enhanced conversion system provides flexible options to handle common compatibility issues between Go and TypeScript protobuf implementations, including:

- Oneof field handling
- Field name transformations
- BigInt serialization
- Default value management

## Key Differences Between Go protojson and protobuf-es

### 1. Go protojson
- Uses `protojson.Marshal()` and `protojson.Unmarshal()`
- Configurable field naming (proto names vs JSON names)
- Oneof fields are serialized as the active field directly
- Handles BigInt as strings automatically
- Can emit or omit default values

### 2. protobuf-es (TypeScript)
- Uses `.toJson()` and `.fromJson()` methods
- Oneof fields may be wrapped in container objects
- BigInt handling requires custom serialization
- Different default value behavior

## Using Conversion Options

### Basic Usage

```typescript
// Create client with default conversion options (auto mode)
const client = new MyServicesClient();

// Or customize with specific converters
const client = new MyServicesClient({
    oneofToJson: 'auto',  // Automatic oneof flattening
    oneofFromJson: 'auto',
    emitDefaults: false,
    bigIntHandler: (value) => value.toString()
});
```

### Custom Oneof Converters

```typescript
// Custom oneof conversion with full control
const client = new MyServicesClient({
    oneofToJson: (oneofValue, context) => {
        // For protobuf-es style: { case: "fieldName", value: {...} }
        if (oneofValue.case && oneofValue.value !== undefined) {
            // Return flattened object for Go
            return { [oneofValue.case]: oneofValue.value };
        }
        return oneofValue;
    },
    oneofFromJson: (jsonValue, context) => {
        // Convert Go's flattened format back to protobuf-es style
        // This would need the schema to know which field was set
        return jsonValue; // Complex without schema
    }
});
```

### Schema-Based Conversion

```typescript
// Provide schema information for better conversion
const schemas: Map<string, MessageSchema> = new Map([
    ['GameMove', {
        name: 'GameMove',
        fields: {
            player: { name: 'player', type: 'scalar' },
            sequenceNum: { name: 'sequenceNum', type: 'scalar' },
            moveType: { 
                name: 'moveType', 
                type: 'oneof',
                oneofFields: ['moveUnit', 'attackUnit', 'buildBase', 'captureBase']
            }
        }
    }]
]);

const client = new MyServicesClient({
    schemaProvider: (messageName) => schemas.get(messageName),
    oneofToJson: 'auto',  // Will use schema information
    oneofFromJson: 'auto'
});
```

### Runtime Configuration

```typescript
// Change conversion options at runtime
client.setConversionOptions({
    fieldTransformer: (field) => {
        // Convert camelCase to snake_case
        return field.replace(/([A-Z])/g, '_$1').toLowerCase();
    }
});
```

## Conversion Options

### `oneofToJson` (function | 'auto')
Converter for oneof fields when sending to WASM. Set to `'auto'` for automatic detection and conversion.

**Auto mode example:**
```typescript
// Input: protobuf-es style
const request = {
    moveType: {
        case: "moveUnit",
        value: { fromQ: -1, fromR: -2, toQ: -1, toR: -1 }
    }
};

// Output: flattened for Go
{
    moveUnit: { fromQ: -1, fromR: -2, toQ: -1, toR: -1 }
}
```

### `oneofFromJson` (function | 'auto')
Converter for oneof fields when receiving from WASM. Set to `'auto'` for automatic handling.

### `schemaProvider` (function)
Optional function that provides message schemas for better conversion accuracy. Returns `MessageSchema` objects with field type information.

### `fieldTransformer` (function)
Custom function to transform field names between TypeScript and Go conventions.

**Example:**
```typescript
{
    fieldTransformer: (field) => {
        // Convert camelCase to snake_case
        if (field === 'userId') return 'user_id';
        if (field === 'createdAt') return 'created_at';
        return field;
    }
}
```

### `emitDefaults` (boolean)
Controls whether default values are included in the JSON.

- `false` (default): Omits zero values, empty strings, empty arrays
- `true`: Includes all fields regardless of value

### `bigIntHandler` (function)
Custom handler for BigInt serialization.

**Example:**
```typescript
{
    bigIntHandler: (value) => {
        // Convert to string for large numbers
        if (value > Number.MAX_SAFE_INTEGER) {
            return value.toString();
        }
        // Use number for smaller values
        return Number(value);
    }
}
```

## WASM-side Configuration

The Go WASM wrapper now uses improved protojson options:

```go
// Unmarshal options
protojson.UnmarshalOptions{
    DiscardUnknown: true,
    AllowPartial:   true,
}

// Marshal options
protojson.MarshalOptions{
    UseProtoNames:   false,  // Use JSON names (camelCase)
    EmitUnpopulated: false,  // Don't emit zero values
    UseEnumNumbers:  false,  // Use enum string values
}
```

## Troubleshooting Common Issues

### 1. Oneof Field Errors
If you're getting errors with oneof fields:

```typescript
// Check your oneof structure
console.log('Request:', JSON.stringify(request, null, 2));

// Ensure auto mode is enabled (default)
client.setConversionOptions({ 
    oneofToJson: 'auto',
    oneofFromJson: 'auto'
});

// Or provide custom converter for complex cases
client.setConversionOptions({
    oneofToJson: (oneofValue, context) => {
        console.log('Converting oneof:', context.fieldName, oneofValue);
        // Custom logic here
        return { [oneofValue.case]: oneofValue.value };
    }
});
```

### 2. Field Name Mismatches
If field names don't match between TypeScript and Go:

```typescript
// Use a field transformer
client.setConversionOptions({
    fieldTransformer: (field) => {
        const transforms = {
            'userId': 'user_id',
            'orderId': 'order_id',
            // Add more mappings as needed
        };
        return transforms[field] || field;
    }
});
```

### 3. BigInt Serialization Issues
For BigInt handling:

```typescript
// Ensure BigInts are properly serialized
client.setConversionOptions({
    bigIntHandler: (value) => value.toString()
});
```

### 4. Missing Required Fields
The WASM side now allows partial messages:

```go
AllowPartial: true  // Won't fail on missing required fields
```

## Best Practices

1. **Test Conversion Early**: Test your proto conversions early in development to identify any compatibility issues.

2. **Use Consistent Naming**: Stick to either camelCase or snake_case consistently across your proto definitions.

3. **Document Oneofs**: Clearly document which fields are oneofs in your proto files to help with debugging.

4. **Handle Errors Gracefully**: The improved error handling provides more context:
   ```typescript
   try {
       const response = await client.service.method(request);
   } catch (error) {
       if (error instanceof WasmError) {
           console.error(`Method: ${error.methodPath}`);
           console.error(`Message: ${error.message}`);
       }
   }
   ```

5. **Monitor Performance**: The conversion system adds minimal overhead, but for high-frequency calls, consider caching conversion options.

## Example: Complete Setup

```typescript
import { MyServicesClient } from './myServicesClient';

// Initialize client with custom options
const client = new MyServicesClient({
    handleOneofs: true,
    emitDefaults: false,
    fieldTransformer: (field) => {
        // Handle specific field mappings
        const mappings: Record<string, string> = {
            'userId': 'user_id',
            'createdAt': 'created_at',
            'updatedAt': 'updated_at'
        };
        return mappings[field] || field;
    },
    bigIntHandler: (value) => value.toString()
});

// Load WASM module
await client.loadWasm('./services.wasm');

// Make a request with a oneof field
const response = await client.userService.searchUsers({
    // This oneof will be properly flattened
    searchBy: {
        email: "user@example.com"
    },
    limit: 10
});
```

## Migration Guide

If you're upgrading from the previous version:

1. The default behavior remains largely the same
2. Oneof handling is now enabled by default
3. You can disable new features by passing `handleOneofs: false`
4. Existing code should continue to work without modifications

## Future Improvements

Planned enhancements include:
- Type-aware field conversion based on proto descriptors
- Automatic detection of oneof fields from type information
- Performance optimizations for large messages
- Better support for Well-Known Types (Timestamp, Duration, etc.)