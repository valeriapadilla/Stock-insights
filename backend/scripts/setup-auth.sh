echo "Setting up authentication for Stock Insights..."
cd backend
go run cmd/setup-auth/main.go

echo ""
echo "Setup complete! You can now:"
echo "1. Start the server: go run cmd/api/main.go"
echo "2. Test the endpoint: curl -X POST http://localhost:8080/api/admin/ingest \\"
echo "   -H 'Authorization: Bearer \$ADMIN_TOKEN'" 