import { useEffect } from 'react';

import { useExhibitionReadinessStore } from '../../stores/exhibitionReadinessStore';
import type { ExhibitionReadinessReason } from '../../types/exhibitionReadiness';

const reasonLabels: Record<string, string> = {
  EMPTY_EXHIBITION: '空展览',
  MISSING_ARTWORK: '作品缺失',
  ARTWORK_DRAFT: '草稿未发布',
  REVIEW_PENDING: '待审核',
  REVIEW_REJECTED: '审核驳回',
  REVIEW_FLAGGED: '审核标记',
  REVIEW_UNKNOWN: '审核状态异常',
  ARTWORK_SOLD: '已售',
  ARTWORK_ARCHIVED: '已下架',
  ARTWORK_STATUS_UNKNOWN: '作品状态异常',
  ARTWORK_NO_IMAGE: '缺少图片',
};

function StatChip({ label, value, tone }: { label: string; value: number; tone: string }) {
  return (
    <div className="flex items-baseline gap-2">
      <span className={`font-display text-2xl ${tone}`}>{value}</span>
      <span className="text-xs text-ink/55">{label}</span>
    </div>
  );
}

function ReasonRow({ reason }: { reason: ExhibitionReadinessReason }) {
  return (
    <li className="border-b border-ink/10 py-3 last:border-b-0">
      <div className="flex items-center justify-between gap-3">
        <span className="inline-flex rounded-full bg-clay/10 px-2.5 py-1 text-xs font-semibold text-clay">
          {reasonLabels[reason.code] ?? reason.code}
        </span>
        <span className="text-xs text-ink/55">{reason.count} 件</span>
      </div>
      <p className="mt-2 text-sm leading-6 text-ink/70">{reason.message}</p>
      {reason.artworkIds.length > 0 && (
        <p className="mt-1 break-all text-xs text-ink/45">作品：{reason.artworkIds.join('、')}</p>
      )}
    </li>
  );
}

export function ExhibitionReadinessPanel({ exhibitionId }: { exhibitionId: string }) {
  const { readinessByExhibition, loadingIds, errorByExhibition, notFoundIds, loadReadiness } =
    useExhibitionReadinessStore();
  const readiness = readinessByExhibition[exhibitionId];
  const loading = loadingIds[exhibitionId] ?? false;
  const error = errorByExhibition[exhibitionId] ?? '';
  const notFound = notFoundIds[exhibitionId] ?? false;

  useEffect(() => {
    void loadReadiness(exhibitionId);
  }, [exhibitionId, loadReadiness]);

  const retry = (
    <button
      onClick={() => void loadReadiness(exhibitionId)}
      disabled={loading}
      className="mt-3 border border-ink px-4 py-2 text-sm hover:bg-ink hover:text-rice disabled:cursor-not-allowed disabled:opacity-50"
    >
      {loading ? '重试中…' : '重新检测'}
    </button>
  );

  if (loading && !readiness) {
    return (
      <section className="border border-ink/15 bg-rice/70 p-6">
        <p className="text-sm text-ink/60">正在实时统计关联作品的就绪情况…</p>
      </section>
    );
  }

  if (error) {
    return (
      <section className="border border-clay/40 bg-clay/5 p-6">
        <h3 className="font-display text-xl text-clay">就绪检测失败</h3>
        <p className="mt-2 text-sm text-ink/70">
          {notFound ? '找不到该展览，就绪统计无法生成（404）。' : `无法获取就绪统计：${error}`}
        </p>
        {!notFound && retry}
      </section>
    );
  }

  if (!readiness) {
    return null;
  }

  const percent = Math.round(readiness.readyRatio * 100);
  const barColor = readiness.ready ? 'bg-moss' : readiness.readyRatio > 0 ? 'bg-lapis' : 'bg-clay';

  return (
    <section className="border border-ink/15 bg-rice/70 p-6">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <p className="text-sm uppercase tracking-[0.2em] text-ink/45">Public readiness</p>
          <h3 className="mt-1 font-display text-2xl">
            {readiness.ready ? '可公开展出' : readiness.isEmpty ? '空展览' : '暂未就绪'}
          </h3>
          <p className="mt-2 max-w-2xl text-sm leading-6 text-ink/70">{readiness.conclusion}</p>
        </div>
        <div className="text-right">
          <p className={`font-display text-4xl ${readiness.ready ? 'text-moss' : 'text-clay'}`}>{percent}%</p>
          <p className="text-xs text-ink/55">
            {readiness.readyCount}/{readiness.total} 件就绪
          </p>
        </div>
      </div>

      <div className="mt-4 h-2 w-full overflow-hidden rounded-full bg-ink/10">
        <div className={`h-full ${barColor} transition-all`} style={{ width: `${percent}%` }} />
      </div>

      <div className="mt-5 flex flex-wrap gap-x-8 gap-y-3">
        <StatChip label="作品总数" value={readiness.total} tone="text-ink" />
        <StatChip label="可公开展出" value={readiness.readyCount} tone="text-moss" />
        <StatChip label="待审核" value={readiness.pendingCount} tone="text-lapis" />
        <StatChip label="已售" value={readiness.soldCount} tone="text-clay" />
        <StatChip label="已下架" value={readiness.archivedCount} tone="text-ink/60" />
        <StatChip label="无图" value={readiness.noImageCount} tone="text-clay" />
      </div>

      {readiness.isEmpty ? (
        <p className="mt-5 rounded-sm border border-dashed border-ink/25 bg-rice p-4 text-sm leading-6 text-ink/65">
          本展尚未收录任何作品，当前就绪比例为 0%。请先在工作台为展览添加作品，再进行公开发布。
        </p>
      ) : readiness.reasons.length > 0 ? (
        <div className="mt-5">
          <p className="text-sm font-semibold text-ink">阻断原因（{readiness.reasons.length}）</p>
          <ul className="mt-2">
            {readiness.reasons.map((reason) => (
              <ReasonRow key={reason.code} reason={reason} />
            ))}
          </ul>
        </div>
      ) : (
        <p className="mt-5 rounded-sm bg-moss/10 p-4 text-sm text-moss">未发现阻断项，所有关联作品均满足公开条件。</p>
      )}

      <div className="mt-5 flex items-center justify-between">
        <span className="text-xs text-ink/45">统计按关联作品实时计算，不使用缓存或模拟数据。</span>
        <button
          onClick={() => void loadReadiness(exhibitionId)}
          disabled={loading}
          className="border border-ink/30 px-4 py-2 text-sm hover:bg-ink hover:text-rice disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading ? '重试中…' : '重新检测'}
        </button>
      </div>
    </section>
  );
}
