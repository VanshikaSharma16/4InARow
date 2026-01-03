# MongoDB Connection Troubleshooting Guide

## Current Error
```
TLS internal error
server selection error: context deadline exceeded
```

This typically means one of these issues:
1. ❌ IP address not whitelisted in MongoDB Atlas
2. ❌ Incorrect connection string format
3. ❌ Password contains special characters that need URL encoding
4. ❌ Network/firewall blocking the connection

## Step-by-Step Fix

### Step 1: Check if .env file exists
```bash
cd /Users/vanshikasharma/Documents/connect4/backend
ls -la .env
```

If it doesn't exist, create it:
```bash
touch .env
```

### Step 2: Whitelist Your IP in MongoDB Atlas

**This is the MOST COMMON fix!**

1. Go to [MongoDB Atlas](https://cloud.mongodb.com/)
2. Log in to your account
3. Click **"Network Access"** in the left sidebar
4. Click **"Add IP Address"**
5. Choose one:
   - **Option A (Easiest for testing)**: Click **"Allow Access from Anywhere"**
     - This adds `0.0.0.0/0` (allows all IPs)
     - ⚠️ Less secure, but works for development
   - **Option B (More secure)**: Click **"Add Current IP Address"**
     - This adds only your current IP
     - You'll need to update if your IP changes
6. Click **"Confirm"**
7. **Wait 1-2 minutes** for the change to propagate

### Step 3: Get Your Connection String

1. In MongoDB Atlas, click **"Clusters"** in the left sidebar
2. Click the **"Connect"** button on your cluster
3. Choose **"Connect your application"**
4. Select **"Go"** (driver: Go, version: latest)
5. Copy the connection string (looks like):
   ```
   mongodb+srv://<username>:<password>@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority
   ```

### Step 4: Create/Update .env File

1. Open/create the `.env` file:
   ```bash
   cd /Users/vanshikasharma/Documents/connect4/backend
   nano .env
   ```

2. Add your connection string, **replacing `<password>` with your actual password**:
   ```
   MONGODB_URI=mongodb+srv://username:YOUR_PASSWORD@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority
   ```

3. **IMPORTANT**: If your password contains special characters, URL-encode them:
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
   - ` ` (space) → `%20`

   **Example**: 
   - Password: `P@ssw0rd#123`
   - Use in connection string: `P%40ssw0rd%23123`

4. Save the file:
   - Press `Ctrl+X`
   - Press `Y` to confirm
   - Press `Enter` to save

### Step 5: Verify Database User

1. In MongoDB Atlas, click **"Database Access"** in the left sidebar
2. Check if your database user exists
3. If needed, create a new user:
   - Click **"Add New Database User"**
   - Choose **"Password"** authentication
   - Username: (choose a username)
   - Password: (create a strong password, write it down!)
   - Database User Privileges: **"Atlas admin"** or **"Read and write to any database"**
   - Click **"Add User"**

### Step 6: Test the Connection

1. Restart your backend server:
   ```bash
   cd /Users/vanshikasharma/Documents/connect4/backend
   go run main.go
   ```

2. You should see:
   ```
   ✅ MongoDB connected successfully
   Server started on :8080
   ```

## Quick Checklist

- [ ] IP address whitelisted in MongoDB Atlas (Network Access)
- [ ] `.env` file exists in `backend` directory
- [ ] Connection string in `.env` has correct format
- [ ] Password in connection string matches MongoDB Atlas user password
- [ ] Special characters in password are URL-encoded
- [ ] Database user exists and has correct permissions
- [ ] Cluster is not paused in MongoDB Atlas

## Alternative: Test Without MongoDB

If you just want to test the game without MongoDB, you can:
1. Delete or rename the `.env` file
2. The server will run without MongoDB (you'll see warnings)
3. Game results won't be saved, but everything else works

## Still Having Issues?

### Try a Simple Password
1. Create a new database user with a simple password (no special characters)
2. Update your `.env` file with the new username and password
3. Test again

### Check Cluster Status
1. Make sure your MongoDB Atlas cluster is **not paused**
2. Check the cluster status in the Atlas dashboard
3. If paused, click "Resume" to start it

### Verify Connection String Format
Your connection string should look like:
```
mongodb+srv://username:password@cluster0.xxxxx.mongodb.net/?retryWrites=true&w=majority
```

NOT:
- ❌ `mongodb://` (use `mongodb+srv://` for Atlas)
- ❌ Missing `?retryWrites=true&w=majority`
- ❌ Incorrect password encoding

### Test with MongoDB Compass
1. Download [MongoDB Compass](https://www.mongodb.com/products/compass)
2. Try connecting with the same connection string
3. If it works in Compass but not in Go, check the Go driver version

