#!/usr/bin/env bash
set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
THREADS="${THREADS:-4}"
CONNECTIONS="${CONNECTIONS:-50}"
DURATION="${DURATION:-5s}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================================="
echo " Starting wrk Benchmark Suite"
echo " Target: $BASE_URL | Threads: $THREADS | Conns: $CONNECTIONS | Duration: $DURATION"
echo "=========================================================="

echo "[1/4] Ensuring test users exist..."
curl -s -X POST "$BASE_URL/students/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"benchstudent1","email":"benchstudent1@example.com","full_name":"Bench Student","password":"Password123!","class":"10A","enrollment_no":"ENR_TEST_001"}' > /dev/null || true

curl -s -X POST "$BASE_URL/teachers/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"benchteacher1","email":"benchteacher1@example.com","full_name":"Bench Teacher","password":"Password123!","class":"10A","subject":"Math"}' > /dev/null || true

curl -s -X POST "$BASE_URL/executives/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"benchexe1","email":"benchexe1@example.com","full_name":"Bench Executive","password":"Password123!","department":"Administration"}' > /dev/null || true

echo "[2/4] Authenticating and obtaining JWT tokens..."
STUDENT_TOKEN=$(curl -s -X POST "$BASE_URL/login" -H "Content-Type: application/json" -d '{"identifier":"benchstudent1","password":"Password123!"}' | grep -oP '(?<="token":")[^"]*')
TEACHER_TOKEN=$(curl -s -X POST "$BASE_URL/login" -H "Content-Type: application/json" -d '{"identifier":"benchteacher1","password":"Password123!"}' | grep -oP '(?<="token":")[^"]*')
EXEC_TOKEN=$(curl -s -X POST "$BASE_URL/login" -H "Content-Type: application/json" -d '{"identifier":"benchexe1","password":"Password123!"}' | grep -oP '(?<="token":")[^"]*')

if [ -z "$STUDENT_TOKEN" ] || [ -z "$TEACHER_TOKEN" ] || [ -z "$EXEC_TOKEN" ]; then
  echo "Error: Failed to obtain authentication tokens. Make sure the API server is running on $BASE_URL."
  exit 1
fi

echo "[3/4] Running Protected Read Endpoints Benchmarks..."
echo ""
echo "--- [1] GET /students/me (Student Profile) ---"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency -H "Authorization: Bearer $STUDENT_TOKEN" "$BASE_URL/students/me"

echo ""
echo "--- [2] GET /teachers/me (Teacher Profile) ---"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency -H "Authorization: Bearer $TEACHER_TOKEN" "$BASE_URL/teachers/me"

echo ""
echo "--- [3] GET /teachers/my-students (Class Students) ---"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency -H "Authorization: Bearer $TEACHER_TOKEN" "$BASE_URL/teachers/my-students"

echo ""
echo "--- [4] GET /executives/me (Executive Profile) ---"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency -H "Authorization: Bearer $EXEC_TOKEN" "$BASE_URL/executives/me"

echo ""
echo "--- [5] GET /executives/students (All Students) ---"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency -H "Authorization: Bearer $EXEC_TOKEN" "$BASE_URL/executives/students"

echo ""
echo "--- [6] GET /executives/teachers (All Teachers) ---"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency -H "Authorization: Bearer $EXEC_TOKEN" "$BASE_URL/executives/teachers"

echo ""
echo "--- [7] GET /executives/all (All Executives) ---"
wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" --latency -H "Authorization: Bearer $EXEC_TOKEN" "$BASE_URL/executives/all"

echo ""
echo "[4/4] Running Authentication Benchmark (Bcrypt)..."
echo "--- [8] POST /login ---"
wrk -t2 -c10 -d"$DURATION" --latency -s "$SCRIPT_DIR/login.lua" "$BASE_URL/login"

echo ""
echo "=========================================================="
echo " Benchmark suite completed successfully!"
echo "=========================================================="
