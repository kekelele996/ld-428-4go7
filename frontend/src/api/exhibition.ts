import { apiPaths } from '../constants/apiPaths';
import { exhibitions } from '../utils/mockData';
import { request } from '../utils/request';
import type { Exhibition, ExhibitionReadiness } from '../types/exhibition';

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

// fetchExhibitionReadiness 实时就绪评估：失败必须抛错，禁止回退模拟数据。
export async function fetchExhibitionReadiness(id: string): Promise<ExhibitionReadiness> {
  return request<ExhibitionReadiness>(apiPaths.exhibitionReadiness(id));
}
