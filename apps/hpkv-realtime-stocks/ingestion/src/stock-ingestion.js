import WebSocket from 'ws';
import axios from 'axios';
import dotenv from 'dotenv';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

// Get the directory path of the current module
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Load environment variables from .env.local
dotenv.config({ path: join(__dirname, '../../.env.local') });

const STOCK_SYMBOLS = ['AAPL', 'GOOGL'];

const HPKV_BASE_URL = process.env.HPKV_BASE_URL;
const HPKV_API_KEY = process.env.HPKV_API_KEY;
const FINNHUB_API_KEY = process.env.FINNHUB_API_KEY;

// Log environment variables (without exposing sensitive data)
console.log('Environment loaded:', {
  HPKV_BASE_URL: HPKV_BASE_URL ? 'Set' : 'Not set',
  HPKV_API_KEY: HPKV_API_KEY ? 'Set' : 'Not set',
  FINNHUB_API_KEY: FINNHUB_API_KEY ? 'Set' : 'Not set'
});

class StockIngestionService {
  constructor() {
    this.connections = new Map();
    this.startAllConnections();
  }

  startAllConnections() {
    STOCK_SYMBOLS.forEach(symbol => {
      this.startConnection(symbol);
    });
  }

  startConnection(symbol) {
    console.log(`Starting connection for ${symbol}`);
    const ws = new WebSocket(`wss://ws.finnhub.io?token=${FINNHUB_API_KEY}`);

    ws.on('open', () => {
      console.log(`Connected to Finnhub WebSocket for ${symbol}`);
      const subscribeMessage = {
        type: 'subscribe',
        symbol: symbol
      };
      console.log(`Sending subscription message for ${symbol}:`, subscribeMessage);
      ws.send(JSON.stringify(subscribeMessage));
    });

    ws.on('message', async (data) => {
      try {
        const message = JSON.parse(data);
        console.log(`Received message for ${symbol}:`, message);
        
        if (message.type === 'trade' && message.data) {
          const latestPrice = message.data[0].p;
          console.log(`Received price for ${symbol}: ${latestPrice}`);
          
          const payload = {
            key: `stock:${symbol}`,
            value: latestPrice.toString()
          };

          // Write the price to HPKV
          await axios.post(
            `${HPKV_BASE_URL}/record`,
            payload,
            {
              headers: {
                'Content-Type': 'application/json',
                'x-api-key': HPKV_API_KEY
              }
            }
          );
          console.log(`Updated HPKV with price for ${symbol}`);
        }
      } catch (error) {
        console.error(`Error processing message for ${symbol}:`, error);
      }
    });

    ws.on('error', (error) => {
      console.error(`WebSocket error for ${symbol}:`, error);
      console.log(error);
      // Attempt to reconnect after error
      setTimeout(() => this.startConnection(symbol), 10000);
    });

    ws.on('close', () => {
      console.log(`WebSocket connection closed for ${symbol}`);
      // Attempt to reconnect after close
      setTimeout(() => this.startConnection(symbol), 10000);
    });

    // Keep the connection alive
    const keepAlive = setInterval(() => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'ping' }));
      }
    }, 30000);

    this.connections.set(symbol, { ws, keepAlive });
  }

  stopConnection(symbol) {
    const connection = this.connections.get(symbol);
    if (connection) {
      clearInterval(connection.keepAlive);
      connection.ws.close();
      this.connections.delete(symbol);
    }
  }

  stopAllConnections() {
    STOCK_SYMBOLS.forEach(symbol => {
      this.stopConnection(symbol);
    });
  }
}

// Create and export a singleton instance
const stockIngestionService = new StockIngestionService();

// Handle process termination
process.on('SIGTERM', () => {
  console.log('SIGTERM received. Stopping all connections...');
  stockIngestionService.stopAllConnections();
  process.exit(0);
});

process.on('SIGINT', () => {
  console.log('SIGINT received. Stopping all connections...');
  stockIngestionService.stopAllConnections();
  process.exit(0);
}); 