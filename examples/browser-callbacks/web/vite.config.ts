import { defineConfig } from 'vite';

export default defineConfig({
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
    }
  },
  
  // Asset handling
  assetsInclude: ['**/*.wasm'],
  
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
