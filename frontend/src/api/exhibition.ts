import { apiPaths, exhibitionReadinessPath } from '../constants/apiPaths';
import { exhibitions } from '../utils/mockData';
import { request } from '../utils/request';
import type { Exhibition } from '../types/exhibition';
import type { ExhibitionReadiness } from '../types/exhibitionReadiness';

export async function fetchExhibitions(): Promise<Exhibition[]> {
  try {
    return await request<Exhibition[]>(apiPaths.exhibitions);
  } catch {
    return exhibitions;
  }
}

export async function fetchExhibition(id: string): Promise<Exhibition | undefined> {
  const list = await fetchExhibitions();
  return list.find((exhibition) => exhibition.id === id);
}

// fetchExhibitionReadiness 实时获取展览就绪统计。
// 注意：接口失败时直接抛错，严禁回退到本地模拟数据。
export async function fetchExhibitionReadiness(id: string): Promise<ExhibitionReadiness> {
  return request<ExhibitionReadiness>(exhibitionReadinessPath(id));
}
