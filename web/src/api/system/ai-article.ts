// 本文件由代码生成器生成，重新生成会覆盖手工修改。
import request from "@/utils/request";

export const aiArticleListApi = (params?: Record<string, unknown>) => request.get("/system/ai-article/index", { params });
export const aiArticleCreateApi = (data: Record<string, unknown>) => request.post("/system/ai-article", data);
export const aiArticleUpdateApi = (id: string | number, data: Record<string, unknown>) => request.put("/system/ai-article/" + id, data);
export const aiArticleDeleteApi = (id: string | number) => request.delete("/system/ai-article/" + id);
