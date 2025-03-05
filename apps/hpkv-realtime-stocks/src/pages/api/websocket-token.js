export default async function handler(req, res) {
  if (req.method !== 'POST') {
    return res.status(405).json({ message: 'Method not allowed' });
  }

  try {
    const { symbol } = req.body;

    // Get WebSocket token from HPKV
    const response = await fetch(`${process.env.HPKV_BASE_URL}/token/websocket`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'x-api-key': process.env.HPKV_API_KEY
      },
      body: JSON.stringify({
        subscribeKeys: [`stock:${symbol}`]
      })
    });

    if (!response.ok) {
      throw new Error('Failed to get WebSocket token');
    }

    const { token } = await response.json();
    return res.status(200).json({ token });
  } catch (error) {
    console.error('Error getting WebSocket token:', error);
    return res.status(500).json({ message: 'Failed to get WebSocket token' });
  }
} 