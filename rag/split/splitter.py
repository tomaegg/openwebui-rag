# coding=utf-8
from fastapi import FastAPI
from pydantic import BaseModel
from langchain.text_splitter import RecursiveCharacterTextSplitter
from typing import List, Optional

app = FastAPI(title="Text Splitter API", version="1.0")


class TextSplitter:
    def __init__(
        self,
        target_chunk_size: int = 512,
        chunk_overlap: int = 53,
        separators: Optional[list[str]] = None,
        keep_separator: bool = False,
    ) -> None:
        """
        文本切分器
        """
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=target_chunk_size,
            chunk_overlap=chunk_overlap,
            separators=separators
            or ["。", "？", "！", "；", "\n", ",", ".", ";", "!", "?"],
            keep_separator=keep_separator,
        )

    def split(self, text: list[str]) -> list[str]:
        """执行文本切分"""
        tmp = self.text_splitter.create_documents(text)
        return [t.page_content for t in tmp]


# ----------- FastAPI 部分 -----------


class SplitRequest(BaseModel):
    text: list[str]
    target_chunk_size: int = 512
    chunk_overlap: int = 53
    separators: Optional[List[str]] = None
    keep_separator: bool = False


class SplitResponse(BaseModel):
    chunks: List[str]
    count: int


splitter = TextSplitter()


@app.post("/split", response_model=SplitResponse)
def split_text(req: SplitRequest):
    chunks = splitter.split(req.text)
    # TODO: 如果不是默认参数，需要新建splitter, 暂时不考虑
    return SplitResponse(chunks=chunks, count=len(chunks))


@app.get("/")
def root():
    return {"message": "Text Splitter API is running. Use POST /split to split text."}
