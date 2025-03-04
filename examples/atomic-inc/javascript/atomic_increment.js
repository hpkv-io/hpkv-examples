import axios from 'axios';
import dotenv from 'dotenv';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';

// Load environment variables from .env file
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
dotenv.config({ path: join(__dirname, '.env') });

// Get HPKV configuration from environment variables
const HPKV_BASE_URL = process.env.HPKV_BASE_URL;
const HPKV_API_KEY = process.env.HPKV_API_KEY;

async function createKey(key, initialValue = 0) {
    try {
        const payload = {
            key,
            value: initialValue.toString()
        };
        console.log(`Creating key with payload: ${JSON.stringify(payload)}`);

        const response = await axios.post(
            `${HPKV_BASE_URL}/record`,
            payload,
            {
                headers: {
                    'Content-Type': 'application/json',
                    'x-api-key': HPKV_API_KEY
                },
                timeout: 10000
            }
        );

        console.log(`Create response status: ${response.status}`);
        console.log(`Create response body: ${JSON.stringify(response.data)}`);

        return true;
    } catch (error) {
        console.error(`Error creating record: ${error.message}`);
        if (error.response) {
            console.error(`Error details: ${JSON.stringify(error.response.data)}`);
        }
        return false;
    }
}

async function atomicIncrement(key, increment) {
    if (!HPKV_BASE_URL || !HPKV_API_KEY) {
        throw new Error('HPKV_BASE_URL and HPKV_API_KEY must be set in environment variables');
    }

    try {
        const payload = {
            key,
            increment
        };
        console.log(`Attempting atomic increment with payload: ${JSON.stringify(payload)}`);

        const response = await axios.post(
            `${HPKV_BASE_URL}/record/atomic`,
            payload,
            {
                headers: {
                    'Content-Type': 'application/json',
                    'x-api-key': HPKV_API_KEY
                }
            }
        );

        console.log(`Atomic increment response status: ${response.status}`);
        console.log(`Atomic increment response body: ${JSON.stringify(response.data)}`);

        // If the key doesn't exist (404), create it first and try again
        if (response.status === 404) {
            console.log(`Key '${key}' doesn't exist. Creating it with initial value 0...`);
            if (!await createKey(key, 0)) {
                throw new Error('Failed to create key with initial value');
            }

            // Retry the increment operation
            console.log('Retrying atomic increment after key creation...');
            const retryResponse = await axios.post(
                `${HPKV_BASE_URL}/record/atomic`,
                payload,
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'x-api-key': HPKV_API_KEY
                    }
                }
            );

            console.log(`Retry response status: ${retryResponse.status}`);
            console.log(`Retry response body: ${JSON.stringify(retryResponse.data)}`);

            if (!retryResponse.data.success) {
                throw new Error(`HPKV API error: ${retryResponse.data.message || 'Unknown error'}`);
            }

            return retryResponse.data;
        }

        if (!response.data.success) {
            throw new Error(`HPKV API error: ${response.data.message || 'Unknown error'}`);
        }

        return response.data;
    } catch (error) {
        if (error.response) {
            throw new Error(`Failed to connect to HPKV API: ${error.message}`);
        }
        throw new Error(`Invalid response from HPKV API: ${error.message}`);
    }
}

async function main() {
    const key = 'counter:example';

    try {
        // Clean up any existing key
        console.log('Cleaning up any existing key...');
        await axios.delete(`${HPKV_BASE_URL}/record/${key}`, {
            headers: {
                'Content-Type': 'application/json',
                'x-api-key': HPKV_API_KEY
            }
        });

        await createKey(key, 0);

        const getResponse = await axios.get(`${HPKV_BASE_URL}/record/${key}`, {
            headers: {
                'Content-Type': 'application/json',
                'x-api-key': HPKV_API_KEY
            }
        });
        console.log(`Get response status: ${getResponse.status}`);

        // Increment by 1
        console.log('\nIncrementing counter by 1...');
        let result = await atomicIncrement(key, 1);
        console.log(`Result: ${JSON.stringify(result)}`);

        // Increment by 5
        console.log('\nIncrementing counter by 5...');
        result = await atomicIncrement(key, 5);
        console.log(`Result: ${JSON.stringify(result)}`);

        // Decrement by 2
        console.log('\nDecrementing counter by 2...');
        result = await atomicIncrement(key, -2);
        console.log(`Result: ${JSON.stringify(result)}`);
    } catch (error) {
        console.error(`Error: ${error.message}`);
        process.exit(1);
    }
}

main(); 