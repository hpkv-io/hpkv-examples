import fetch from 'node-fetch';

export class HPKVClient {
    constructor(baseUrl = null, apiKey = null) {
        this.baseUrl = (baseUrl || process.env.HPKV_BASE_URL || '').replace(/\/$/, '');
        this.apiKey = apiKey || process.env.HPKV_API_KEY;

        if (!this.baseUrl) {
            throw new Error('HPKV base URL not provided. Set HPKV_BASE_URL environment variable or pass baseUrl parameter.');
        }
        if (!this.apiKey) {
            throw new Error('HPKV API key not provided. Set HPKV_API_KEY environment variable or pass apiKey parameter.');
        }

        this.headers = {
            'Content-Type': 'application/json',
            'x-api-key': this.apiKey
        };
    }

    _serializeValue(value) {
        if (typeof value === 'string') {
            return value;
        }
        return JSON.stringify(value);
    }

    async create(key, value) {
        try {
            const payload = {
                key,
                value: this._serializeValue(value)
            };

            console.log(`Sending request to ${this.baseUrl}/record`);
            console.log('Payload:', JSON.stringify(payload, null, 2));

            const response = await fetch(`${this.baseUrl}/record`, {
                method: 'POST',
                headers: this.headers,
                body: JSON.stringify(payload)
            });

            if (!response.ok) {
                const error = await response.text();
                console.error(`Create failed with status ${response.status} - ${error}`);
                return false;
            }

            console.log(`Create succeeded with status ${response.status}`);
            return true;
        } catch (error) {
            console.error('Error creating record:', error.message);
            return false;
        }
    }

    async read(key) {
        try {
            const response = await fetch(`${this.baseUrl}/record/${key}`, {
                headers: this.headers
            });

            if (!response.ok) {
                const error = await response.text();
                console.error(`Error in read: Status ${response.status} - ${error}`);
                return null;
            }

            const data = await response.json();
            if (!data || !data.value) {
                return null;
            }

            try {
                return JSON.parse(data.value);
            } catch {
                return data.value;
            }
        } catch (error) {
            console.error('Error reading record:', error.message);
            return null;
        }
    }

    async update(key, value, partialUpdate = false) {
        try {
            const payload = {
                key,
                value: this._serializeValue(value),
                partialUpdate
            };

            const response = await fetch(`${this.baseUrl}/record`, {
                method: 'POST',
                headers: this.headers,
                body: JSON.stringify(payload)
            });

            if (!response.ok) {
                const error = await response.text();
                console.error(`Update failed with status ${response.status} - ${error}`);
                return false;
            }

            return true;
        } catch (error) {
            console.error('Error updating record:', error.message);
            return false;
        }
    }

    async delete(key) {
        try {
            const response = await fetch(`${this.baseUrl}/record/${key}`, {
                method: 'DELETE',
                headers: this.headers
            });

            if (!response.ok) {
                const error = await response.text();
                console.error(`Delete failed with status ${response.status} - ${error}`);
                return false;
            }

            return true;
        } catch (error) {
            console.error('Error deleting record:', error.message);
            return false;
        }
    }
} 