import request from "@/utils/request";

export const aiArticleListApi = (params?: any) => request.get("/system/ai-article/index", { params });
export const aiArticleCreateApi = (data: any) => request.post("/system/ai-article", data);
export const aiArticleUpdateApi = (id: string | number, data: any) => request.put("/system/ai-article/" + id, data);
export const aiArticleDeleteApi = (id: string | number) => request.delete("/system/ai-article/" + id);
