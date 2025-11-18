#!/bin/bash

# Development script to run both frontend and backend in development mode

echo "Starting Weather Agency Label Tool in development mode..."
echo ""

# Check if tmux is available
if command -v tmux &> /dev/null; then
    echo "Using tmux to run both services..."
    
    # Create a new tmux session
    SESSION="weather-label-dev"
    
    # Kill existing session if it exists
    tmux kill-session -t $SESSION 2>/dev/null
    
    # Create new session with backend
    tmux new-session -d -s $SESSION -n backend "cd $(pwd) && echo 'Starting backend...' && go run main.go"
    
    # Create new window for frontend
    tmux new-window -t $SESSION -n frontend "cd $(pwd)/frontend && echo 'Starting frontend...' && npm run dev"
    
    # Attach to the session
    echo "Services started in tmux session '$SESSION'"
    echo "Press Ctrl+B then D to detach"
    echo "Run 'tmux attach -t $SESSION' to reattach"
    echo ""
    tmux attach -t $SESSION
else
    echo "tmux not found. Starting services sequentially..."
    echo ""
    echo "Starting frontend development server..."
    echo "Backend will need to be started separately with: go run main.go"
    echo ""
    cd frontend && npm run dev
fi
