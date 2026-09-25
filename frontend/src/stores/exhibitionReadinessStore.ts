import { create } from 'zustand';

import { fetchExhibitionReadiness } from '../api/exhibition';
import { ApiError } from '../utils/request';
import type { ExhibitionReadiness } from '../types/exhibitionReadiness';

interface ExhibitionReadinessState {
  readinessByExhibition: Record<string, ExhibitionReadiness>;
  loadingIds: Record<string, boolean>;
  errorByExhibition: Record<string, string>;
  notFoundIds: Record<string, boolean>;
  loadReadiness: (exhibitionId: string) => Promise<void>;
}

export const useExhibitionReadinessStore = create<ExhibitionReadinessState>((set) => ({
  readinessByExhibition: {},
  loadingIds: {},
  errorByExhibition: {},
  notFoundIds: {},
  loadReadiness: async (exhibitionId) => {
    set((state) => ({
      loadingIds: { ...state.loadingIds, [exhibitionId]: true },
      errorByExhibition: { ...state.errorByExhibition, [exhibitionId]: '' },
      notFoundIds: { ...state.notFoundIds, [exhibitionId]: false },
    }));
    try {
      // 只使用后端实时数据，接口失败不回退模拟数据。
      const readiness = await fetchExhibitionReadiness(exhibitionId);
      set((state) => ({
        readinessByExhibition: { ...state.readinessByExhibition, [exhibitionId]: readiness },
        loadingIds: { ...state.loadingIds, [exhibitionId]: false },
      }));
    } catch (err) {
      const message = err instanceof Error ? err.message : '就绪统计接口请求失败';
      set((state) => ({
        loadingIds: { ...state.loadingIds, [exhibitionId]: false },
        errorByExhibition: { ...state.errorByExhibition, [exhibitionId]: message },
        notFoundIds: {
          ...state.notFoundIds,
          [exhibitionId]: err instanceof ApiError && err.status === 404,
        },
      }));
    }
  },
}));
