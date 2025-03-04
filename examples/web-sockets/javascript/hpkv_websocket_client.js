const WebSocket = require('ws');

const OperationCode = Object.freeze({
    GET: 1,
    INSERT: 2,
    UPDATE: 3,
    DELETE: 4
});

class HPKVWebSocketClient {
    constructor(baseUrl, apiKey) {
        this.baseUrl = baseUrl.replace(/^http/, 'ws');
        this.apiKey = apiKey;
        this.messageId = 0;
        this.responseFutures = new Map();
        this.ws = null;
        this.connect();
    }

    connect() {
        return new Promise((resolve, reject) => {
            try {
                this.ws = new WebSocket(`${this.baseUrl}/ws?apiKey=${this.apiKey}`);
                
                this.ws.on('open', () => {
                    console.log('WebSocket connection established');
                    resolve();
                });

                this.ws.on('message', (data) => {
                    try {
                        const response = JSON.parse(data);
                        const msgId = response.messageId;
                        if (msgId && this.responseFutures.has(msgId)) {
                            const { resolve, reject } = this.responseFutures.get(msgId);
                            this.responseFutures.delete(msgId);
                            
                            if (response.error) {
                                reject(new Error(response.error));
                            } else {
                                resolve(response);
                            }
                        }
                    } catch (error) {
                        console.error('Error handling message:', error);
                    }
                });

                this.ws.on('close', (code, reason) => {
                    console.log(`WebSocket connection closed: ${code} - ${reason}`);
                });

                this.ws.on('error', (error) => {
                    console.error('WebSocket error:', error);
                    reject(error);
                });
            } catch (error) {
                reject(error);
            }
        });
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    async sendMessage(message) {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            await this.connect();
        }

        return new Promise((resolve, reject) => {
            try {
                const msgId = ++this.messageId;
                message.messageId = msgId;
                
                this.responseFutures.set(msgId, { resolve, reject });
                this.ws.send(JSON.stringify(message));
            } catch (error) {
                reject(error);
            }
        });
    }

    async create(key, value) {
        try {
            const message = {
                op: OperationCode.INSERT,
                key,
                value: typeof value === 'string' ? value : JSON.stringify(value)
            };
            
            const response = await this.sendMessage(message);
            return !response.error;
        } catch (error) {
            console.error('Error creating record:', error);
            return false;
        }
    }

    async read(key) {
        try {
            const message = {
                op: OperationCode.GET,
                key
            };
            
            const response = await this.sendMessage(message);
            if (response.error) {
                return null;
            }
            
            try {
                return JSON.parse(response.value);
            } catch {
                return response.value;
            }
        } catch (error) {
            console.error('Error reading record:', error);
            return null;
        }
    }

    async update(key, value, partialUpdate = false) {
        try {
            const message = {
                op: partialUpdate ? OperationCode.UPDATE : OperationCode.INSERT,
                key,
                value: typeof value === 'string' ? value : JSON.stringify(value)
            };
            
            const response = await this.sendMessage(message);
            return !response.error;
        } catch (error) {
            console.error('Error updating record:', error);
            return false;
        }
    }

    async delete(key) {
        try {
            const message = {
                op: OperationCode.DELETE,
                key
            };
            
            const response = await this.sendMessage(message);
            return !response.error;
        } catch (error) {
            console.error('Error deleting record:', error);
            return false;
        }
    }
}

module.exports = HPKVWebSocketClient; 