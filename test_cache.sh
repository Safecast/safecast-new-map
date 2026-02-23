#!/bin/bash

# Test script to verify the cache functionality works

echo "Starting the improved Safecast server..."

# Start the server in the background
./safecast-new-map-improved -port 8766 &
SERVER_PID=$!

# Wait a moment for the server to start
sleep 3

# Check if the server is running
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "Server is running with PID $SERVER_PID"
    
    # Test the cache admin endpoint
    echo "Testing cache admin endpoints..."
    
    # Test cache stats
    echo "Cache stats:"
    curl -X POST "http://localhost:8766/api/admin/cache?action=stats&password=" 2>/dev/null || echo "Could not reach cache stats endpoint"
    
    # Test cache clear
    echo "Clearing cache:"
    curl -X POST "http://localhost:8766/api/admin/cache?action=clear&password=" 2>/dev/null || echo "Could not reach cache clear endpoint"
    
    echo "Cache functionality test completed."
else
    echo "Failed to start server"
fi

# Kill the server
if kill -0 $SERVER_PID 2>/dev/null; then
    kill $SERVER_PID
    echo "Server stopped."
fi