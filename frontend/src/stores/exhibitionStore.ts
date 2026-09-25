import { create } from 'zustand';

import { fetchExhibitionReadiness, fetchExhibitions } from '../api/exhibition';
import type { Exhibition, ExhibitionReadiness } from '../types/exhibition';

interface ExhibitionState {
  exhibitions: Exhibition[];
  readiness: ExhibitionReadiness | null;
  readinessLoading: boolean;
  readinessError: string | null;
  loadExhibitions: () => Promise<void>;
  loadReadiness: (id: string) => Promise<void>;
}

export const useExhibitionStore = create<ExhibitionState>((set) => ({
  exhibitions: [],
  readiness: null,
  readinessLoading: false,
  readinessError: null,
  loadExhibitions: async () => set({ exhibitions: await fetchExhibitions() }),
  // 就绪评估失败时记录错误并清空旧数据，由页面展示重试入口，绝不回退模拟数据。
  loadReadiness: async (id) => {
    set({ readinessLoading: true, readinessError: null });
    try {
      const readiness = await fetchExhibitionReadiness(id);
      set({ readiness, readinessLoading: false });
    } catch (err) {
      set({
        readiness: null,
        readinessLoading: false,
        readinessError: err instanceof Error ? err.message : '就绪情况加载失败',
      });
    }
  },
}));
