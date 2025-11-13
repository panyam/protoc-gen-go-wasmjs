// Import the generated base bundle and service clients
import { ExampleBundle } from './generated';
import { PresenterServiceClient } from './generated/presenter/v1/presenterServiceClient';
import { BrowserAPIClient } from './generated/browser/v1/browserAPIClient';

// Import TypeScript types for better type safety
import type { 
  LoadUserDataRequest, 
  StateUpdateRequest, 
  PreferencesRequest,
  CallbackDemoRequest 
} from './generated/presenter/v1/interfaces';



// Types for better code organization
interface BrowserAPIImpl {
  fetch(request: any): Promise<any>;
  getLocalStorage(request: any): Promise<any>;
  setLocalStorage(request: any): Promise<any>;
  getCookie(request: any): Promise<any>;
  alert(request: any): Promise<any>;
  promptUser(request: any): Promise<any>;
  logToWindow(request: any): Promise<any>;
}

// Log function
function log(message: string, type: 'info' | 'error' | 'success' = 'info') {
  const output = document.getElementById('output');
  if (!output) return;
  
  const timestamp = new Date().toLocaleTimeString();
  const prefix = type === 'error' ? '❌' : type === 'success' ? '✅' : 'ℹ️';
  output.textContent += `[${timestamp}] ${prefix} ${message}\n`;
  output.scrollTop = output.scrollHeight;
  console.log(`[${type}]`, message);
}

// Add UI update to list
function addUIUpdate(update: any) {
  const list = document.getElementById('uiUpdates');
  if (!list) return;
  
  const item = document.createElement('li');
  item.className = 'ui-update';
  item.innerHTML = `
    <strong>${update.component}.${update.action}</strong>
    <pre>${JSON.stringify(update.data, null, 2)}</pre>
  `;
  if (list.firstChild) {
    list.insertBefore(item, list.firstChild);
  } else {
    list.appendChild(item);
  }

  // Keep only last 5 updates
  while (list.children.length > 5) {
    const lastChild = list.lastChild;
    if (lastChild) {
      list.removeChild(lastChild);
    }
  }
}

// Update status
function setStatus(message: string, type: 'info' | 'loading' | 'success' | 'error' = 'info') {
  const status = document.getElementById('status');
  if (!status) return;
  
  status.textContent = message;
  status.className = `status ${type}`;
}

// Browser API Implementation
class BrowserAPIImpl implements BrowserAPIImpl {
  async fetch(request: any) {
    log(`Fetch called: ${request.method} ${request.url}`);

    // Simulate API response for demo
    if (request.url.includes('/users/')) {
      return {
        status: 200,
        statusText: 'OK',
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify({
          username: 'john_doe',
          email: 'john@example.com',
          permissions: ['read', 'write', 'admin']
        })
      };
    }

    return {
      status: 404,
      statusText: 'Not Found',
      headers: {},
      body: ''
    };
  }

  async getLocalStorage(request: any) {
    log(`GetLocalStorage: ${request.key}`);
    const value = localStorage.getItem(request.key);
    return {
      value: value || '',
      exists: value !== null
    };
  }

  async setLocalStorage(request: any) {
    log(`SetLocalStorage: ${request.key} = ${request.value}`);
    try {
      localStorage.setItem(request.key, request.value);
      return { success: true };
    } catch (e: any) {
      log(`Failed to set localStorage: ${e.message}`, 'error');
      return { success: false };
    }
  }

  async getCookie(request: any) {
    log(`GetCookie: ${request.name}`);
    const cookies = document.cookie.split(';');
    for (const cookie of cookies) {
      const [name, value] = cookie.trim().split('=');
      if (name === request.name) {
        return { value, exists: true };
      }
    }
    return { value: '', exists: false };
  }

  async alert(request: any) {
    log(`Alert: ${request.message}`);
    alert(request.message);
    return { shown: true };
  }

  async promptUser(request: any) {
    log(`PromptUser: ${request.message}`);
    // Note: protobuf uses snake_case (default_value) but JavaScript typically uses camelCase
    const defaultValue = request.defaultValue || request.default_value || '';
    const result = window.prompt(request.message, defaultValue);
    return {
      value: result || '',
      cancelled: result === null
    };
  }

  async logToWindow(request: any) {
    const logOutput = document.getElementById('log-output');
    if (logOutput) {
      const logEntry = document.createElement('div');
      const timestamp = new Date().toLocaleTimeString();

      const levelColors = {
        'error': '#ff0000',
        'warning': '#ff9900',
        'success': '#00cc00',
        'info': '#0066cc'
      };

      const color = levelColors[request.level as keyof typeof levelColors] || '#000000';
      logEntry.style.color = color;
      logEntry.textContent = `[${timestamp}] [${request.level || 'info'}] ${request.message}`;

      logOutput.appendChild(logEntry);
      logOutput.scrollTop = logOutput.scrollHeight;
    }
    log(`LogToWindow [${request.level}]: ${request.message}`);
    return { logged: true };
  }
}

// Initialize the application
async function init() {
  try {
    setStatus('Loading WASM module...', 'loading');

    // Create base bundle with module configuration
    const wasmBundle = new ExampleBundle();

    // Create service clients using composition
    const presenterService = new PresenterServiceClient(wasmBundle);
    const browserAPI = new BrowserAPIClient(wasmBundle);

    // Register browser API implementation
    wasmBundle.registerBrowserService('BrowserAPI', new BrowserAPIImpl());

    // Load WASM module
    await wasmBundle.loadWasm('/browser_example.wasm');

    setStatus('WASM module loaded successfully!', 'success');
    log('WASM module initialized', 'success');

    // Enable buttons
    const buttons = ['loadUserBtn', 'updateStateBtn', 'savePrefsBtn', 'callbackDemoBtn'];
    buttons.forEach(id => {
      const btn = document.getElementById(id) as HTMLButtonElement;
      if (btn) btn.disabled = false;
    });

    // Wire up button handlers
    setupEventHandlers(presenterService);
  } catch (error: any) {
    setStatus(`Failed to initialize: ${error.message}`, 'error');
    log(`Initialization failed: ${error.message}`, 'error');
    console.error(error);
  }
}

function setupEventHandlers(presenterService: PresenterServiceClient) {
  // Load User Data button
  const loadUserBtn = document.getElementById('loadUserBtn');
  loadUserBtn?.addEventListener('click', async () => {
    const userIdInput = document.getElementById('userId') as HTMLInputElement;
    const userId = userIdInput?.value || 'user123';
    log(`Loading user data for: ${userId}`);

    try {
      const request: LoadUserDataRequest = {
        userId: userId
      };
      const response = await presenterService.loadUserData(request);

      log(`User loaded: ${response.username} (${response.email})`, 'success');
      log(`Permissions: ${response.permissions.join(', ')}`);
      log(`From cache: ${response.fromCache}`);
    } catch (error: any) {
      log(`Failed to load user: ${error.message}`, 'error');
    }
  });

  // Update State button
  const updateStateBtn = document.getElementById('updateStateBtn');
  updateStateBtn?.addEventListener('click', async () => {
    const actionSelect = document.getElementById('stateAction') as HTMLSelectElement;
    const action = actionSelect?.value || 'refresh';
    log(`Updating UI state: ${action}`);

    const params: Record<string, string> = action === 'navigate' ?
      { page: '/dashboard' } :
      { timestamp: new Date().toISOString() };

    try {
      const request: StateUpdateRequest = { action, params };
      await presenterService.updateUIState(
        request,
        (update, error, done) => {
          if (error) {
            log(`Stream error: ${error}`, 'error');
            return false;
          }
          if (done) {
            log('UI update stream complete', 'success');
            return false;
          }

          addUIUpdate(update);
          return true; // Continue stream
        }
      );
    } catch (error: any) {
      log(`Failed to update state: ${error.message}`, 'error');
    }
  });

  // Save Preferences button
  const savePrefsBtn = document.getElementById('savePrefsBtn');
  savePrefsBtn?.addEventListener('click', async () => {
    const preferences: Record<string, string> = {};

    const prefKey1 = (document.getElementById('prefKey1') as HTMLInputElement)?.value;
    const prefValue1 = (document.getElementById('prefValue1') as HTMLInputElement)?.value;
    if (prefKey1) preferences[prefKey1] = prefValue1;

    const prefKey2 = (document.getElementById('prefKey2') as HTMLInputElement)?.value;
    const prefValue2 = (document.getElementById('prefValue2') as HTMLInputElement)?.value;
    if (prefKey2) preferences[prefKey2] = prefValue2;

    log(`Saving preferences: ${JSON.stringify(preferences)}`);

    try {
      const request: PreferencesRequest = {
        preferences
      };
      const response = await presenterService.savePreferences(request);

      log(`Preferences saved: ${response.itemsSaved} items`, 'success');
    } catch (error: any) {
      log(`Failed to save preferences: ${error.message}`, 'error');
    }
  });

  // Callback Demo button
  const callbackDemoBtn = document.getElementById('callbackDemoBtn') as HTMLButtonElement;
  callbackDemoBtn?.addEventListener('click', async () => {
    const resultDiv = document.getElementById('callbackResult');

    callbackDemoBtn.disabled = true;
    callbackDemoBtn.textContent = '⏳ Running demo...';
    if (resultDiv) {
      resultDiv.innerHTML = '<div style="color: #856404;">Demo in progress... Watch for prompts!</div>';
    }

    // Clear previous logs
    const logOutput = document.getElementById('log-output');
    if (logOutput) logOutput.innerHTML = '';

    try {
      const request: CallbackDemoRequest = {
        demoName: 'User Input Collection'
      };
      await presenterService.runCallbackDemo(request, (response, error) => {
        if (error) {
          throw new Error(error);
        }
        
        // Handle the response when the async method completes
        handleCallbackDemoResponse(response);
      });
    } catch (error: any) {
      if (resultDiv) {
        resultDiv.innerHTML = `<div style="color: #721c24;">Error: ${error.message}</div>`;
      }
      log(`Callback demo error: ${error.message}`, 'error');
      
      // Re-enable button on error
      callbackDemoBtn.disabled = false;
      callbackDemoBtn.textContent = 'Start Callback Demo';
    }
  });
}

// Handle callback demo response
function handleCallbackDemoResponse(response: any) {
  const callbackDemoBtn = document.getElementById('callbackDemoBtn') as HTMLButtonElement;
  const resultDiv = document.getElementById('callbackResult');

  if (resultDiv) {
    if (response.completed) {
      resultDiv.innerHTML = `
        <div style="color: #155724; background: #d4edda; padding: 10px; border-radius: 4px;">
          ✅ Demo completed!<br>
          Collected: ${response.collectedInputs.join(', ')}
        </div>
      `;
    } else {
      resultDiv.innerHTML = `
        <div style="color: #721c24; background: #f8d7da; padding: 10px; border-radius: 4px;">
          ❌ Demo was cancelled<br>
          Partial: ${response.collectedInputs.join(', ') || 'None'}
        </div>
      `;
    }
  }

  // Re-enable button
  if (callbackDemoBtn) {
    callbackDemoBtn.disabled = false;
    callbackDemoBtn.textContent = 'Start Callback Demo';
  }
}

// Start initialization when page loads
document.addEventListener('DOMContentLoaded', init);
