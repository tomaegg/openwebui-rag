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
        chunk_overlap: int = 0,
        separators: Optional[list[str]] = None,
        keep_separator: bool = False,
    ) -> None:
        """
        文本切分器
        """
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=target_chunk_size,
            chunk_overlap=chunk_overlap,
            separators=separators or ["。", "？", "！", "；", "\n"],
            keep_separator=keep_separator,
        )

    def split(self, text: str) -> list[str]:
        """执行文本切分"""
        tmp = self.text_splitter.create_documents([text])
        return [t.page_content for t in tmp]


# ----------- FastAPI 部分 -----------

class SplitRequest(BaseModel):
    text: str
    target_chunk_size: int = 512
    chunk_overlap: int = 0
    separators: Optional[List[str]] = None
    keep_separator: bool = False


class SplitResponse(BaseModel):
    chunks: List[str]
    count: int


@app.post("/split", response_model=SplitResponse)
def split_text(req: SplitRequest):
    splitter = TextSplitter(
        target_chunk_size=req.target_chunk_size,
        chunk_overlap=req.chunk_overlap,
        separators=req.separators,
        keep_separator=req.keep_separator,
    )
    chunks = splitter.split(req.text)
    return SplitResponse(chunks=chunks, count=len(chunks))


@app.get("/")
def root():
    return {"message": "Text Splitter API is running. Use POST /split to split text."}
