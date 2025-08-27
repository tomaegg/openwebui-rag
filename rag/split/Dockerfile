FROM python:3.11-slim

# 安装 uv
RUN pip install --no-cache-dir uv

WORKDIR /app

# 拷贝 pyproject.toml 和 uv.lock
COPY pyproject.toml uv.lock ./

# 用 uv 安装依赖
RUN uv sync --frozen

# 拷贝代码
COPY . .

EXPOSE 8000

CMD ["uv", "run", "--", "uvicorn", "splitter:app", "--host", "0.0.0.0", "--port", "8000"]
