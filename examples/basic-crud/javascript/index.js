import 'dotenv/config';
import { HPKVClient } from './hpkv-client.js';

async function main() {
    try {
        // Initialize HPKV client using environment variables
        const client = new HPKVClient();

        console.log('HPKV CRUD Operations Example');
        console.log('===========================');
        console.log(`\nUsing HPKV server: ${client.baseUrl}`);

        // Create operation
        const userData = {
            name: 'John Doe',
            email: 'john@example.com',
            age: 30
        };

        console.log('\n1. Creating a new user record...');
        const success = await client.create('user:1', userData);
        if (!success) {
            console.log('Failed to create record. Exiting...');
            process.exit(1);
        }
        console.log('Create operation succeeded');

        // Read operation
        console.log('\n2. Reading the user record...');
        const retrievedData = await client.read('user:1');
        if (retrievedData) {
            console.log('Retrieved data:', JSON.stringify(retrievedData, null, 2));
        } else {
            console.log('Failed to retrieve data');
        }

        // Update operation
        console.log('\n3. Updating the user\'s age...');
        userData.age = 31;
        const updateSuccess = await client.update('user:1', userData);
        console.log(`Update operation ${updateSuccess ? 'succeeded' : 'failed'}`);

        // Read after update
        console.log('\n4. Reading the updated user record...');
        const updatedData = await client.read('user:1');
        if (updatedData) {
            console.log('Retrieved data:', JSON.stringify(updatedData, null, 2));
        } else {
            console.log('Failed to retrieve data');
        }

        // Delete operation
        console.log('\n5. Deleting the user record...');
        const deleteSuccess = await client.delete('user:1');
        console.log(`Delete operation ${deleteSuccess ? 'succeeded' : 'failed'}`);

        // Verify deletion
        console.log('\n6. Attempting to read deleted record...');
        const deletedData = await client.read('user:1');
        if (deletedData === null) {
            console.log('Record was successfully deleted');
        } else {
            console.log('Record still exists');
        }

    } catch (error) {
        console.error('Error running example:', error.message);
        process.exit(1);
    }
}

main(); 