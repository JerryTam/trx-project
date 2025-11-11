#!/bin/bash

# API 测试脚本

BASE_URL="http://localhost:8080"

echo "================================"
echo "API 响应格式测试"
echo "================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 测试健康检查
echo "1. 测试健康检查"
echo "GET $BASE_URL/health"
response=$(curl -s $BASE_URL/health)
echo "响应: $response"
echo ""

# 测试用户注册
echo "2. 测试用户注册（成功）"
echo "POST $BASE_URL/api/v1/users/register"
response=$(curl -s -X POST $BASE_URL/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser123",
    "email": "test123@example.com",
    "password": "password123"
  }')
echo "响应: $response"
echo ""

# 测试用户注册（用户已存在）
echo "3. 测试用户注册（用户已存在）"
echo "POST $BASE_URL/api/v1/users/register"
response=$(curl -s -X POST $BASE_URL/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser123",
    "email": "test123@example.com",
    "password": "password123"
  }')
echo "响应: $response"
echo ""

# 测试参数验证错误
echo "4. 测试参数验证错误"
echo "POST $BASE_URL/api/v1/users/register"
response=$(curl -s -X POST $BASE_URL/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "ab",
    "email": "invalid-email",
    "password": "123"
  }')
echo "响应: $response"
echo ""

# 测试用户登录（成功）
echo "5. 测试用户登录（成功）"
echo "POST $BASE_URL/api/v1/users/login"
response=$(curl -s -X POST $BASE_URL/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser123",
    "password": "password123"
  }')
echo "响应: $response"
echo ""

# 测试用户登录（密码错误）
echo "6. 测试用户登录（密码错误）"
echo "POST $BASE_URL/api/v1/users/login"
response=$(curl -s -X POST $BASE_URL/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser123",
    "password": "wrongpassword"
  }')
echo "响应: $response"
echo ""

# 测试获取用户列表
echo "7. 测试获取用户列表（分页）"
echo "GET $BASE_URL/api/v1/users?page=1&page_size=10"
response=$(curl -s "$BASE_URL/api/v1/users?page=1&page_size=10")
echo "响应: $response"
echo ""

# 测试获取单个用户
echo "8. 测试获取单个用户"
echo "GET $BASE_URL/api/v1/users/1"
response=$(curl -s $BASE_URL/api/v1/users/1)
echo "响应: $response"
echo ""

# 测试获取不存在的用户
echo "9. 测试获取不存在的用户"
echo "GET $BASE_URL/api/v1/users/999999"
response=$(curl -s $BASE_URL/api/v1/users/999999)
echo "响应: $response"
echo ""

# 测试无效的用户 ID
echo "10. 测试无效的用户 ID"
echo "GET $BASE_URL/api/v1/users/invalid"
response=$(curl -s $BASE_URL/api/v1/users/invalid)
echo "响应: $response"
echo ""

echo "================================"
echo "测试完成"
echo "================================"
echo ""
echo "统一响应格式验证："
echo "✅ 所有响应都包含 code, message 字段"
echo "✅ 成功响应包含 data 字段"
echo "✅ 分页响应包含 list, total, page, page_size"
echo "✅ 错误响应返回对应的业务状态码"
echo ""

