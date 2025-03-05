# Real-time Stock Dashboard with HPKV

This is a sample application demonstrating HPKV's real-time pub/sub capabilities using a stock price dashboard. The application shows real-time stock price updates for selected companies using HPKV's WebSocket feature and Finnhub's stock data API.

## Features

- Real-time stock price updates using HPKV's pub/sub system
- Interactive line chart showing price changes over time
- WebSocket-based updates without manual refreshing
- Secure API key management


## Prerequisites

- Node.js 18 or later
- HPKV API endpoint and API key ([Sign up for free](https://hpkv.io/signup))
- Finnhub API key ([Get your free API key](https://finnhub.io/dashboard))

## HPKV Pub/Sub Feature

HPKV provides a powerful real-time pub/sub system that enables instant data synchronization across your applications. In this example, we use HPKV's WebSocket API to:

1. Subscribe to real-time stock price updates
2. Receive instant notifications when prices change
3. Maintain persistent connections for live data streaming

The pub/sub system is particularly useful for:
- Real-time data synchronization
- Live updates without polling
- Scalable event-driven architectures
- Instant notifications across multiple clients

For more details about HPKV's WebSocket API, visit the [official documentation](https://hpkv.io/docs/websocket-api).

## Environment Setup

1. Create a `.env.local` file in the root directory with the following variables:

```env
# HPKV Configuration
HPKV_API_KEY=your_hpkv_api_key        # Get from https://hpkv.io/dashboard/api-keys
HPKV_BASE_URL=https://api.hpkv.io
NEXT_PUBLIC_HPKV_BASE_URL=https://api.hpkv.io

# Finnhub Configuration
FINNHUB_API_KEY=your_finnhub_api_key  # Get from https://finnhub.io/dashboard
```

## Installation

1. Install dependencies and start the website:
```bash
npm install
npm run dev
```

2. Start the ingestion service (in a separate terminal):
```bash
cd ingestion
npm install
npm start
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser.

Note: The ingestion service needs to be running for the dashboard to receive real-time stock updates. Make sure to keep both the ingestion service and the development server running.

## Architecture

The application is built with a modern architecture that separates concerns and ensures efficient data flow:

### 1. Data Ingestion Service
Located in `ingestion/src/stock-ingestion.js`, this service:
- Connects to Finnhub's WebSocket API
- Receives real-time stock price updates
- Writes data to HPKV using the key format `stock:SYMBOL`
- Runs independently of the frontend application

### 2. Frontend Application
The Next.js frontend (`src/pages/index.js`):
- Connects to HPKV's WebSocket endpoint
- Subscribes to HPKV keys for stock price updates
- Displays real-time data in an interactive chart
- Handles user interactions and stock selection

### 3. API Routes
- `/api/websocket-token`: Securely generates WebSocket tokens for frontend connections
- Handles API key management and security

## Key Code Sections

### WebSocket Connection
At the client side, the method below is used to start a websocket connection to HPKV and subscribe to the selected keys
```javascript

// Connect to HPKV WebSocket
    const connectWebSocket = async () => {
      try {
        setError(null);
        
        // Get WebSocket token
        const response = await fetch('/api/websocket-token', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ symbol: selectedStock.symbol }),
        });

        if (!response.ok) {
          throw new Error('Failed to get WebSocket token');
        }

        const { token } = await response.json();

        // Create WebSocket connection
        const websocket = new WebSocket(`${process.env.NEXT_PUBLIC_HPKV_BASE_URL}/ws?token=${token}`);

        websocket.onopen = () => {
          console.log('Connected to HPKV WebSocket');
          setError(null);
        };

        websocket.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            
            if (data.type === 'notification') {
              const newData = {
                timestamp: new Date(data.timestamp).toLocaleTimeString(),
                price: parseFloat(data.value)
              };
              setStockData(prev => [...prev, newData].slice(-50)); // Keep last 50 data points
            }
          } catch (error) {
            console.error('Error processing WebSocket message:', error);
          }
        };

        websocket.onerror = (error) => {
          console.error('WebSocket error:', error);
          setError('Connection error. Please try again.');
        };

        websocket.onclose = () => {
          console.log('WebSocket connection closed');
        };

        setWs(websocket);

        return () => {
          if (websocket.readyState === WebSocket.OPEN) {
            websocket.close();
          }
        };
      } catch (error) {
        console.error('Error connecting to WebSocket:', error);
        setError('Failed to connect to WebSocket. Please try again.');
      }
    };
```

### Data Ingestion
Ingestion service opens a websocket to Funhub API to get the stock updates and writes them to HPKV databse.
```javascript
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
```

## Available Stocks

- AAPL (Apple Inc.)
- GOOGL (Alphabet Inc.)


## Technologies Used

- Next.js for the frontend application
- React for UI components
- Recharts for interactive charts
- HPKV WebSocket API for real-time updates
- Finnhub WebSocket API for stock data
- TailwindCSS for styling

## Security Considerations

- API keys are stored securely in environment variables
- WebSocket tokens are generated server-side
- Frontend never exposes sensitive credentials
- Secure WebSocket connections (wss://)

## License

MIT
