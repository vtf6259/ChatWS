import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react'; // Import React plugin

export default defineConfig({
  plugins: [
    react() // Include the React plugin for JSX support
  ],
  server: {
    host: 'localhost', // The server host
    port: 3000,       // The port number
    strictPort: true, // Ensure the server fails if the port is in use
    hmr: {
      // HMR options (optional)
      protocol: 'ws',  // WebSocket protocol for hot module replacement
      host: 'localhost', // HMR host (typically the same as server host)
      port: 3000,      // HMR port (can match the server port)
    },
  },
});
