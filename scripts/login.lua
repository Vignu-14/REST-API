-- wrk script for benchmarking POST /login
wrk.method = "POST"
wrk.body   = '{"identifier":"benchexe1","password":"Password123!"}'
wrk.headers["Content-Type"] = "application/json"
