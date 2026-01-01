# MongoDB Atlas Setup Guide

## Current Status
- ✅ Connection string is being read correctly
- ✅ Username: `mongo`
- ✅ Password: mongo12345
- ❌ Authentication failing

## Step-by-Step Fix

### Step 1: Verify/Create Database User

1. Go to [MongoDB Atlas](https://cloud.mongodb.com/)
2. Click **"Database Access"** in the left sidebar
3. Check if a user named `mongo` exists:
   - **If it exists**: Click "Edit" → Verify the password matches what's in your `.env` file
   - **If it doesn't exist**: Click "Add New Database User"
     - Choose "Password" authentication
     - Username: `mongo` (or any username you prefer)
     - Password: Create a strong password (write it down!)
     - Database User Privileges: "Atlas admin" (for full access) or "Read and write to any database"
     - Click "Add User"

### Step 2: Whitelist Your IP Address

1. In MongoDB Atlas, click **"Network Access"** in the left sidebar
2. Click **"Add IP Address"**
3. You have two options:
   - **Option A (Recommended for testing)**: Click "Allow Access from Anywhere" 
     - This adds `0.0.0.0/0` (allows all IPs)
     - ⚠️ Less secure, but good for development
   - **Option B (More secure)**: Click "Add Current IP Address"
     - This adds only your current IP
     - You'll need to update this if your IP changes

### Step 3: Get Fresh Connection String

1. In MongoDB Atlas, go to your **cluster** (click "Clusters" in the left sidebar)
2. Click the **"Connect"** button on your cluster
3. Choose **"Connect your application"**
4. Select **"Go"** (driver: Go, version: latest)
5. Copy the connection string (it will look like):
   ```
   mongodb+srv://<username>:<password>@cluster0.i22q9fg.mongodb.net/?retryWrites=true&w=majority
   ```

### Step 4: Update Your .env File

1. Open your `.env` file in the `backend` directory:
   ```bash
   cd /Users/vanshikasharma/Documents/connect4/backend
   nano .env
   ```

2. Replace the connection string with the one from Step 3, but **replace `<password>` with your actual password**:
   ```
   MONGODB_URI=mongodb+srv://mongo:YOUR_ACTUAL_PASSWORD@cluster0.i22q9fg.mongodb.net/?retryWrites=true&w=majority
   ```

3. **Important**: If your password contains special characters, they need to be URL-encoded:
   - `@` → `%40`
   - `#` → `%23`
   - `$` → `%24`
   - `%` → `%25`
   - `&` → `%26`
   - `/` → `%2F`
   - `:` → `%3A`
   - `?` → `%3F`
   - `+` → `%2B`
   - `=` → `%3D`

   **Example**: If your password is `P@ssw0rd#123`, use `P%40ssw0rd%23123` in the connection string.

4. Save the file (Ctrl+X, then Y, then Enter)

### Step 5: Test the Connection

Run the test script:
```bash
go run test_connection.go
```

Or run your main application:
```bash
go run main.go
```

## Quick Fix Checklist

- [ ] Database user `mongo` exists in MongoDB Atlas
- [ ] Password in `.env` matches the password in MongoDB Atlas
- [ ] Your IP address is whitelisted in Network Access (or `0.0.0.0/0` is added)
- [ ] Connection string in `.env` has the correct format
- [ ] Password special characters are URL-encoded (if any)

## Still Having Issues?

If authentication still fails after following all steps:

1. **Reset the database user password**:
   - Go to Database Access → Edit user → Change password
   - Update your `.env` file with the new password

2. **Try creating a new database user**:
   - Create a user with a simple password (no special characters)
   - Update your connection string with the new username and password

3. **Verify the cluster is running**:
   - Make sure your MongoDB Atlas cluster is not paused
   - Check the cluster status in the Atlas dashboard


