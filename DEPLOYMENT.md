# Deployment Guide

This guide will help you deploy your Connect 4 game to production. We'll deploy the frontend and backend separately.

## 📋 Prerequisites

- GitHub account (for code hosting)
- MongoDB Atlas account (already set up)
- Accounts on deployment platforms (see options below)

## 🚀 Deployment Options

### Frontend Deployment (React App)

#### Option 1: Vercel (Recommended - Easiest)

1. **Install Vercel CLI:**
   ```bash
   npm install -g vercel
   ```

2. **Build and deploy:**
   ```bash
   cd frontend
   npm run build
   vercel
   ```

3. **Or use Vercel Dashboard:**
   - Go to [vercel.com](https://vercel.com)
   - Sign up/login with GitHub
   - Click "New Project"
   - Import your GitHub repository
   - Set root directory to `frontend`
   - Build command: `npm run build`
   - Output directory: `build`
   - Add environment variable: `REACT_APP_WS_URL=wss://your-backend-url.com/ws`
   - Deploy!

**Note:** Update `WS_URL` in `App.js` to use your backend URL instead of `localhost:8080`

#### Option 2: Netlify

1. **Install Netlify CLI:**
   ```bash
   npm install -g netlify-cli
   ```

2. **Build and deploy:**
   ```bash
   cd frontend
   npm run build
   netlify deploy --prod
   ```

3. **Or use Netlify Dashboard:**
   - Go to [netlify.com](https://netlify.com)
   - Sign up/login with GitHub
   - Click "New site from Git"
   - Connect your repository
   - Set build command: `cd frontend && npm install && npm run build`
   - Set publish directory: `frontend/build`
   - Add environment variable: `REACT_APP_WS_URL=wss://your-backend-url.com/ws`
   - Deploy!

#### Option 3: GitHub Pages

1. **Install gh-pages:**
   ```bash
   cd frontend
   npm install --save-dev gh-pages
   ```

2. **Update package.json:**
   ```json
   {
     "homepage": "https://yourusername.github.io/connect4",
     "scripts": {
       "predeploy": "npm run build",
       "deploy": "gh-pages -d build"
     }
   }
   ```

3. **Deploy:**
   ```bash
   npm run deploy
   ```

### Backend Deployment (Go Server)

#### Option 1: Railway (Recommended - Easy & Free tier)

1. **Go to [railway.app](https://railway.app)**
2. **Sign up/login with GitHub**
3. **Click "New Project" → "Deploy from GitHub repo"**
4. **Select your repository**
5. **Add a new service → Select your backend folder**
6. **Configure:**
   - Root directory: `backend`
   - Build command: `go build -o server`
   - Start command: `./server`
   - Add environment variable: `MONGODB_URI=your-mongodb-connection-string`
   - Add environment variable: `PORT=8080` (Railway will set this automatically)
7. **Deploy!**

**Note:** Railway will give you a URL like `https://your-app.railway.app`

#### Option 2: Render

1. **Go to [render.com](https://render.com)**
2. **Sign up/login with GitHub**
3. **Click "New" → "Web Service"**
4. **Connect your repository**
5. **Configure:**
   - Name: `connect4-backend`
   - Environment: `Go`
   - Root directory: `backend`
   - Build command: `go build -o server`
   - Start command: `./server`
   - Add environment variables:
     - `MONGODB_URI=your-mongodb-connection-string`
     - `PORT=8080`
6. **Deploy!**

#### Option 3: Fly.io

1. **Install Fly CLI:**
   ```bash
   curl -L https://fly.io/install.sh | sh
   ```

2. **Login:**
   ```bash
   fly auth login
   ```

3. **Initialize:**
   ```bash
   cd backend
   fly launch
   ```

4. **Configure fly.toml:**
   ```toml
   [env]
     MONGODB_URI = "your-mongodb-connection-string"
     PORT = "8080"
   ```

5. **Deploy:**
   ```bash
   fly deploy
   ```

#### Option 4: DigitalOcean App Platform

1. **Go to [digitalocean.com](https://digitalocean.com)**
2. **Create account**
3. **Click "Create" → "App"**
4. **Connect GitHub repository**
5. **Configure:**
   - Add component: Backend service
   - Source directory: `backend`
   - Build command: `go build -o server`
   - Run command: `./server`
   - Environment variables:
     - `MONGODB_URI=your-mongodb-connection-string`
     - `PORT=8080`
6. **Deploy!**

## 🔧 Configuration Steps

### Step 1: Update Frontend WebSocket URL

After deploying backend, update the frontend to use the production WebSocket URL:

**In `frontend/src/App.js`:**
```javascript
// Change from:
const WS_URL = "ws://localhost:8080/ws";

// To (for production):
const WS_URL = process.env.REACT_APP_WS_URL || "wss://your-backend-url.com/ws";
```

**Or use environment variable:**
- Create `.env` file in `frontend` directory:
  ```
  REACT_APP_WS_URL=wss://your-backend-url.com/ws
  ```

### Step 2: Update Backend CORS (if needed)

**In `backend/main.go`:**
```go
func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Update with your frontend URL
		w.Header().Set("Access-Control-Allow-Origin", "https://your-frontend-url.com")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// ... rest of the code
	}
}
```

### Step 3: Update MongoDB Atlas Network Access

1. Go to MongoDB Atlas dashboard
2. Click "Network Access"
3. Add your deployment platform's IP ranges or use `0.0.0.0/0` for testing
4. Save changes

### Step 4: Update Leaderboard API URL

**In `frontend/src/Leaderboard.js`:**
```javascript
// Change from:
fetch("http://localhost:8080/leaderboard")

// To:
const API_URL = process.env.REACT_APP_API_URL || "https://your-backend-url.com";
fetch(`${API_URL}/leaderboard`)
```

## 📝 Complete Deployment Checklist

### Backend:
- [ ] Code pushed to GitHub
- [ ] MongoDB Atlas connection string ready
- [ ] Backend deployed to chosen platform
- [ ] Backend URL obtained
- [ ] Environment variables set (MONGODB_URI, PORT)
- [ ] Backend is accessible and running

### Frontend:
- [ ] WebSocket URL updated to production backend URL
- [ ] API URL updated for leaderboard
- [ ] Environment variables set (if using)
- [ ] Frontend deployed to chosen platform
- [ ] Frontend URL obtained
- [ ] Tested connection to backend

### MongoDB:
- [ ] Network Access configured for deployment platform
- [ ] Connection string tested
- [ ] Database user has correct permissions

### Testing:
- [ ] Frontend connects to backend WebSocket
- [ ] Can start a new game
- [ ] Bot joins after 10 seconds
- [ ] Can make moves
- [ ] Leaderboard loads and updates
- [ ] Reconnection works

## 🌐 Example Deployment URLs

After deployment, you'll have URLs like:

- **Frontend:** `https://connect4-game.vercel.app`
- **Backend:** `https://connect4-backend.railway.app`
- **WebSocket:** `wss://connect4-backend.railway.app/ws`
- **API:** `https://connect4-backend.railway.app/leaderboard`

## 🐛 Common Deployment Issues

### WebSocket Connection Failed
- **Issue:** Frontend can't connect to backend WebSocket
- **Fix:** 
  - Use `wss://` (secure WebSocket) for HTTPS sites
  - Check CORS settings
  - Verify backend URL is correct

### CORS Errors
- **Issue:** Browser blocks API requests
- **Fix:** Update CORS headers in backend to allow your frontend domain

### MongoDB Connection Failed
- **Issue:** Backend can't connect to MongoDB
- **Fix:**
  - Add deployment platform IP to MongoDB Atlas Network Access
  - Verify connection string is correct
  - Check environment variables are set

### Build Failures
- **Issue:** Frontend/Backend won't build
- **Fix:**
  - Check build logs for errors
  - Ensure all dependencies are in package.json/go.mod
  - Verify build commands are correct

## 💡 Quick Start (Recommended: Vercel + Railway)

1. **Deploy Backend to Railway:**
   - Sign up at railway.app
   - Connect GitHub repo
   - Deploy backend folder
   - Copy backend URL

2. **Deploy Frontend to Vercel:**
   - Sign up at vercel.com
   - Connect GitHub repo
   - Set root to `frontend`
   - Add env var: `REACT_APP_WS_URL=wss://your-railway-url.up.railway.app/ws`
   - Deploy

3. **Update MongoDB Atlas:**
   - Add Railway IP ranges to Network Access
   - Or use `0.0.0.0/0` for testing

4. **Test:**
   - Open frontend URL
   - Start a game
   - Verify everything works!

## 📚 Additional Resources

- [Vercel Documentation](https://vercel.com/docs)
- [Railway Documentation](https://docs.railway.app)
- [Render Documentation](https://render.com/docs)
- [MongoDB Atlas Documentation](https://docs.atlas.mongodb.com)

---

**Good luck with your deployment! 🚀**

