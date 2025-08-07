// Enhanced client that integrates with the new factory/deserializer system
// This demonstrates how to use the schema-aware deserializer with cross-package factory composition

import { LibraryV2Factory } from './gen/wasm/all-services/library/v2/factory';
import { LibraryV2Deserializer } from './gen/wasm/all-services/library/v2/library_deserializer';
import { LibraryV2SchemaRegistry } from './gen/wasm/all-services/library/v2/library_schemas';
import { Book, FindBooksRequest, FindBooksResponse } from './gen/wasm/all-services/library/v2/library_interfaces';

/**
 * Enhanced client that uses the schema-aware deserializer system
 */
export class EnhancedLibraryClient {
    private factory: LibraryV2Factory;
    private deserializer: LibraryV2Deserializer;
    private baseClient: any; // Reference to the original WASM client

    constructor(baseClient: any) {
        this.baseClient = baseClient;
        this.factory = new LibraryV2Factory();
        this.deserializer = new LibraryV2Deserializer(LibraryV2SchemaRegistry, this.factory);
    }

    /**
     * Enhanced findBooks method with proper deserialization
     */
    async findBooks(request: FindBooksRequest): Promise<FindBooksResponse> {
        // Call the original WASM client method
        const rawResponse = await this.baseClient.libraryService.findBooks(request);
        
        // Use the enhanced deserializer to create properly typed response
        const response = this.deserializer.createAndDeserialize<FindBooksResponse>(
            'library.v2.FindBooksResponse',
            rawResponse
        );
        
        if (!response) {
            throw new Error('Failed to deserialize FindBooksResponse');
        }
        
        return response;
    }

    /**
     * Create a new Book instance using the factory
     */
    createBook(data?: any): Book {
        const result = this.factory.newBook(undefined, undefined, undefined, data);
        if (data && !result.fullyLoaded) {
            // Use deserializer to populate the instance
            return this.deserializer.deserialize(result.instance, data, 'library.v2.Book');
        }
        return result.instance;
    }

    /**
     * Create a new FindBooksRequest instance using the factory
     */
    createFindBooksRequest(data?: any): FindBooksRequest {
        const result = this.factory.newFindBooksRequest(undefined, undefined, undefined, data);
        if (data && !result.fullyLoaded) {
            // Use deserializer to populate the instance
            return this.deserializer.deserialize(result.instance, data, 'library.v2.FindBooksRequest');
        }
        return result.instance;
    }

    /**
     * Demonstrate cross-package factory composition
     * This creates objects that contain BaseMessage from the common package
     */
    createBookWithBase(data?: any): Book {
        const book = this.createBook(data);
        
        // If data contains base field, it will be properly created using the commonFactory
        // through the getFactoryMethod delegation system
        
        return book;
    }

    /**
     * Test factory composition with complex nested data
     */
    async testComplexDeserialization(): Promise<void> {
        // Example data with cross-package references
        const testData = {
            metadata: {
                request_id: "test-123",
                user_agent: "enhanced-client/1.0",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer token123"
                }
            },
            books: [
                {
                    base: {
                        id: "book-1",
                        timestamp: Date.now(),
                        version: "1.0"
                    },
                    title: "Advanced TypeScript",
                    author: "Jane Doe",
                    isbn: "978-1234567890",
                    year: 2023,
                    genre: "Technology",
                    available: true,
                    tags: ["typescript", "programming", "web"],
                    rating: 4.8
                },
                {
                    base: {
                        id: "book-2",
                        timestamp: Date.now(),
                        version: "1.0"
                    },
                    title: "Modern JavaScript",
                    author: "John Smith",
                    isbn: "978-0987654321",
                    year: 2023,
                    genre: "Technology",
                    available: false,
                    tags: ["javascript", "es6", "nodejs"],
                    rating: 4.5
                }
            ],
            totalCount: 2,
            hasMore: false
        };

        console.log('Testing complex deserialization with cross-package dependencies...');
        
        // Create and deserialize a FindBooksResponse
        const response = this.deserializer.createAndDeserialize<FindBooksResponse>(
            'library.v2.FindBooksResponse',
            testData
        );

        if (response) {
            console.log('✅ Successfully deserialized FindBooksResponse');
            console.log('📊 Response metadata:', response.metadata);
            console.log('📚 Number of books:', response.books?.length);
            
            if (response.books && response.books.length > 0) {
                const firstBook = response.books[0];
                console.log('📖 First book:', firstBook.title);
                console.log('🏷️ First book base:', firstBook.base);
                console.log('🔗 Cross-package composition working:', !!firstBook.base);
            }
        } else {
            console.error('❌ Failed to deserialize response');
        }
    }
}

// Example usage
export async function demonstrateEnhancedClient() {
    // Import the original WASM client
    const { default: OriginalClient } = await import('./gen/wasm/all-services/library_all_servicesClient.client');
    
    // Create instances
    const originalClient = new OriginalClient();
    const enhancedClient = new EnhancedLibraryClient(originalClient);
    
    // Wait for WASM to load (in a real app)
    // await originalClient.loadWasm();
    // await originalClient.waitUntilReady();
    
    // Test the enhanced functionality
    await enhancedClient.testComplexDeserialization();
    
    console.log('🎉 Enhanced client demonstration complete!');
}