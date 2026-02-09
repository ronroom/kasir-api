#!/bin/bash
BASE_URL="http://localhost:8080"

echo "Running Verification..."

# Check if server is reachable
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo "Error: Server is not running at $BASE_URL. Please make sure the API is running (e.g., in another terminal: go run main.go)."
    exit 1
fi

# 1. Test Negative Price (Product)
echo "1. Testing Negative Price..."
RESPONSE=$(curl -s -X POST "$BASE_URL/api/produk" -d '{"name":"Test Bad Price","price":-5000,"stock":10,"category_id":1}' -H "Content-Type: application/json")
if [[ "$RESPONSE" == *"Price cannot be negative"* ]]; then
  echo "PASS: Negative Price Rejected"
else
  echo "FAIL: Negative Price Accepted or Wrong Error: $RESPONSE"
fi

# 2. Test Negative Stock (Product)
echo "2. Testing Negative Stock..."
RESPONSE=$(curl -s -X POST "$BASE_URL/api/produk" -d '{"name":"Test Bad Stock","price":5000,"stock":-10,"category_id":1}' -H "Content-Type: application/json")
if [[ "$RESPONSE" == *"Stock cannot be negative"* ]]; then
  echo "PASS: Negative Stock Rejected"
else
  echo "FAIL: Negative Stock Accepted or Wrong Error: $RESPONSE"
fi

# 3. Create Valid Product for Checkout Tests
echo "3. Creating Valid Product..."
VALID_PRODUCT=$(curl -s -X POST "$BASE_URL/api/produk" -d '{"name":"Test Checkout Item","price":1000,"stock":5,"category_id":1}' -H "Content-Type: application/json")
# Extract ID using grep if available, or just rely on manual check if complicated json parsing is needed without jq
# Assuming simple regex works for ID
PRODUCT_ID=$(echo $VALID_PRODUCT | grep -o '"id":[0-9]*' | grep -o '[0-9]*')

if [ -z "$PRODUCT_ID" ]; then
    echo "FAIL: Could not create product. Response: $VALID_PRODUCT"
    exit 1
fi
echo "Created Product ID: $PRODUCT_ID"

# 4. Test Checkout Negative Quantity
echo "4. Testing Checkout Negative Quantity..."
RESPONSE=$(curl -s -X POST "$BASE_URL/api/checkout" -d "{\"items\":[{\"product_id\":$PRODUCT_ID,\"quantity\":-1}]}" -H "Content-Type: application/json")
if [[ "$RESPONSE" == *"Quantity must be greater than 0"* ]]; then
  echo "PASS: Negative Quantity Rejected"
else
  echo "FAIL: Negative Quantity Accepted or Wrong Error: $RESPONSE"
fi

# 5. Test Checkout Insufficient Stock
echo "5. Testing Checkout Insufficient Stock (Req: 10, Stock: 5)..."
RESPONSE=$(curl -s -X POST "$BASE_URL/api/checkout" -d "{\"items\":[{\"product_id\":$PRODUCT_ID,\"quantity\":10}]}" -H "Content-Type: application/json")
if [[ "$RESPONSE" == *"insufficient stock"* ]]; then
  echo "PASS: Insufficient Stock Rejected"
else
  echo "FAIL: Insufficient Stock Allowed or Wrong Error: $RESPONSE"
fi

# 6. Test Successful Checkout
echo "6. Testing Successful Checkout (Req: 2, Stock: 5)..."
RESPONSE=$(curl -s -X POST "$BASE_URL/api/checkout" -d "{\"items\":[{\"product_id\":$PRODUCT_ID,\"quantity\":2}]}" -H "Content-Type: application/json")
TRANSACTION_ID=$(echo $RESPONSE | grep -o '"id":[0-9]*' | grep -o '[0-9]*')

if [ -n "$TRANSACTION_ID" ]; then
    echo "PASS: Checkout Successful. Transaction ID: $TRANSACTION_ID"
else
    echo "FAIL: Checkout Failed. Response: $RESPONSE"
    echo "FAIL: Checkout Failed. Response body: $RESPONSE"
fi

# 7. Test Stock Update (Should be 3)
echo "7. Verifying Stock Update..."
PRODUCT_INFO=$(curl -s "$BASE_URL/api/produk/$PRODUCT_ID")
STOCK=$(echo $PRODUCT_INFO | grep -o '"stock":[0-9]*' | grep -o '[0-9]*')
if [ "$STOCK" == "3" ]; then
    echo "PASS: Stock Updated Correctly (3)"
else
    echo "FAIL: Stock Incorrect ($STOCK). Expected 3."
fi

# Cleanup
echo "Cleaning up..."
curl -s -X DELETE "$BASE_URL/api/produk/$PRODUCT_ID" > /dev/null

# 8. Test Daily Report
echo "8. Testing Daily Report..."
REPORT=$(curl -s "$BASE_URL/api/report/hari-ini")
TOTAL_TRANS=$(echo $REPORT | grep -o '"total_transaksi":[0-9]*' | grep -o '[0-9]*')
if [ -n "$TOTAL_TRANS" ] && [ "$TOTAL_TRANS" -gt 0 ]; then
    echo "PASS: Daily Report returned transactions: $TOTAL_TRANS"
else
    echo "FAIL: Daily Report failed or empty. Response: $REPORT"
fi

# 9. Test Date Range Report (Challenge)
echo "9. Testing Date Range Report..."
# Use a wide range to ensure we catch today's transactions
TODAY=$(date +%Y-%m-%d)
RANGE_REPORT=$(curl -s "$BASE_URL/api/report?start_date=$TODAY&end_date=$TODAY")
RANGE_TRANS=$(echo $RANGE_REPORT | grep -o '"total_transaksi":[0-9]*' | grep -o '[0-9]*')
if [ "$RANGE_TRANS" == "$TOTAL_TRANS" ]; then
    echo "PASS: Date Range Report matches Daily Report ($RANGE_TRANS)"
else
    echo "FAIL: Date Range Report mismatch. Expected $TOTAL_TRANS, got $RANGE_TRANS. Response: $RANGE_REPORT"
fi

echo "Verification Complete."
