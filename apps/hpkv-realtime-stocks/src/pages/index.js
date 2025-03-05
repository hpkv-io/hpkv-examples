import { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ReferenceLine } from 'recharts';

const STOCK_SYMBOLS = [
  { symbol: 'AAPL', name: 'Apple Inc.' },
  { symbol: 'GOOGL', name: 'Alphabet Inc.' },
];

export default function Dashboard() {
  const [selectedStock, setSelectedStock] = useState(STOCK_SYMBOLS[0]);
  const [stockData, setStockData] = useState([]);
  const [error, setError] = useState(null);
  const [zoomLevel, setZoomLevel] = useState(1);
  const [yAxisDomain, setYAxisDomain] = useState(['auto', 'auto']);

  // Calculate Y-axis domain based on data
  const calculateYAxisDomain = (data) => {
    if (data.length === 0) return ['auto', 'auto'];
    
    const prices = data.map(d => d.price);
    const minPrice = Math.min(...prices);
    const maxPrice = Math.max(...prices);
    const priceRange = maxPrice - minPrice;
    
    // Add padding to the range (10% on each side)
    const padding = priceRange * 0.1;
    return [minPrice - padding, maxPrice + padding];
  };

  useEffect(() => {
    if (stockData.length > 0) {
      setYAxisDomain(calculateYAxisDomain(stockData));
    }
  }, [stockData]);

  useEffect(() => {
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

    connectWebSocket();
  }, [selectedStock]);

  const handleResetZoom = () => {
    setZoomLevel(1);
    setYAxisDomain(['auto', 'auto']);
  };

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <h1 className="text-3xl font-bold text-white mb-8">Real-time Stock Dashboard</h1>
          
          {/* Stock Selection */}
          <div className="mb-8">
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Select Stock
            </label>
            <select
              value={selectedStock.symbol}
              onChange={(e) => {
                const stock = STOCK_SYMBOLS.find(s => s.symbol === e.target.value);
                setSelectedStock(stock);
                setStockData([]);
                setError(null);
                handleResetZoom();
              }}
              className="hpkv-select mt-1 block w-full pl-3 pr-10 py-2 text-base focus:outline-none sm:text-sm"
            >
              {STOCK_SYMBOLS.map((stock) => (
                <option key={stock.symbol} value={stock.symbol}>
                  {stock.name} ({stock.symbol})
                </option>
              ))}
            </select>
          </div>

          {/* Error Message */}
          {error && (
            <div className="mb-4 p-4 bg-red-900/50 border border-red-700 rounded-lg">
              <p className="text-red-300">{error}</p>
            </div>
          )}

          {/* Stock Chart */}
          <div className="hpkv-card p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold text-white">
                {selectedStock.name} ({selectedStock.symbol}) Price Chart
              </h2>
              
            </div>
            <div className="h-[400px]">
              {stockData.length > 0 ? (
                <LineChart
                  width={800}
                  height={400}
                  data={stockData}
                  margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                  scale={zoomLevel}
                >
                  <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                  <XAxis 
                    dataKey="timestamp" 
                    stroke="#9CA3AF"
                    tick={{ fill: '#9CA3AF' }}
                  />
                  <YAxis 
                    stroke="#9CA3AF"
                    tick={{ fill: '#9CA3AF' }}
                    domain={yAxisDomain}
                    tickFormatter={(value) => `$${value.toFixed(2)}`}
                  />
                  <Tooltip 
                    contentStyle={{ 
                      backgroundColor: '#1F2937', 
                      border: '1px solid #374151',
                      borderRadius: '0.5rem',
                      color: '#F3F4F6'
                    }}
                    labelStyle={{ color: '#9CA3AF' }}
                    formatter={(value) => [`$${value.toFixed(2)}`, 'Price']}
                  />
                  <Legend 
                    wrapperStyle={{ color: '#9CA3AF' }}
                  />
                  <ReferenceLine
                    y={stockData[stockData.length - 1]?.price}
                    stroke="#8B5CF6"
                    strokeDasharray="3 3"
                    label={{
                      value: `Current: $${stockData[stockData.length - 1]?.price.toFixed(2)}`,
                      position: 'right',
                      fill: '#9CA3AF'
                    }}
                  />
                  <Line
                    type="monotone"
                    dataKey="price"
                    stroke="#8B5CF6"
                    strokeWidth={2}
                    dot={false}
                  />
                </LineChart>
              ) : (
                <div className="flex items-center justify-center h-full">
                  <p className="text-gray-400">Waiting for data...</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
