import { defineConfig } from 'vite';

export default defineConfig({
  // Configure WASM file handling
  define: {
    global: 'globalThis',
  },
  // Root directory for the project
  root: '.',
  
  // Public directory for static assets
  publicDir: 'public',
  
  // Build configuration
  build: {
    outDir: 'dist',
    sourcemap: true,
    target: 'es2020',
    rollupOptions: {
      input: {
        main: './index.html'
      }
    }
  },
  
  // Dev server configuration
  server: {
    port: 8080,
    open: true,
    cors: true,
    // Serve WASM files with correct MIME type
    headers: {
      'Cross-Origin-Embedder-Policy': 'require-corp',
      'Cross-Origin-Opener-Policy': 'same-origin'
    },
    // Configure MIME types for WASM files
    middlewareMode: false,
    fs: {
      strict: false
    }
  },
  
  // Asset handling
  assetsInclude: ['**/*.wasm'],
  
  // Plugin configuration for WASM
  plugins: [],
  
  // Path resolution
  resolve: {
    alias: {
      '@/gen': '../gen/wasm/ts'
    }
  },
  
  // Optimize dependencies
  optimizeDeps: {
    include: ['@protoc-gen-go-wasmjs/runtime']
  }
});
