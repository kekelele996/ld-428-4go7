import type { ExhibitionReadiness } from '../../types/exhibition';
import { EmptyState } from './EmptyState';
import { StatCard } from './StatCard';

interface ExhibitionReadinessPanelProps {
  readiness: ExhibitionReadiness | null;
  loading: boolean;
  error: string | null;
  onRetry: () => void;
}

// ExhibitionReadinessPanel 展览公开就绪评估面板：展示结论、实时统计、阻断原因与重试入口。
export function ExhibitionReadinessPanel({ readiness, loading, error, onRetry }: ExhibitionReadinessPanelProps) {
  return (
    <section className="border border-ink/15 bg-rice p-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <h3 className="font-display text-2xl text-ink">公开就绪评估</h3>
        <button
          type="button"
          onClick={onRetry}
          disabled={loading}
          className="border border-ink/25 px-4 py-2 text-sm text-ink transition hover:border-clay hover:text-clay disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading ? '评估中…' : '重新评估'}
        </button>
      </div>

      {loading && <p className="mt-6 text-sm text-ink/60">正在实时统计关联作品…</p>}

      {!loading && error && (
        <div className="mt-6 border border-dashed border-clay/50 bg-clay/5 p-6 text-center">
          <p className="text-sm text-clay">就绪情况加载失败：{error}</p>
          <button
            type="button"
            onClick={onRetry}
            className="mt-4 border border-clay px-4 py-2 text-sm text-clay transition hover:bg-clay hover:text-rice"
          >
            重试
          </button>
        </div>
      )}

      {!loading && !error && readiness && (
        <div className="mt-6">
          <p className={`inline-block px-3 py-1 text-sm ${readiness.ready ? 'bg-moss/15 text-moss' : 'bg-clay/15 text-clay'}`}>
            {readiness.ready ? '结论：已就绪，全部作品均可公开展出' : '结论：未就绪，存在阻断公开的原因'}
          </p>

          {readiness.totalArtworks === 0 ? (
            <div className="mt-6">
              <EmptyState
                title="空展览"
                description="该展览尚未收录任何作品，暂无可公开内容。请先在创作工作台为展览添加作品后再重新评估。"
              />
            </div>
          ) : (
            <>
              <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
                <StatCard label="关联作品总数" value={readiness.totalArtworks} />
                <StatCard label="可公开展出" value={readiness.publicReadyCount} />
                <StatCard label="待审核" value={readiness.pendingReviewCount} />
                <StatCard label="已售" value={readiness.soldCount} />
                <StatCard label="已下架" value={readiness.archivedCount} />
                <StatCard label="无图" value={readiness.noImageCount} />
                <StatCard label="就绪比例" value={`${Math.round(readiness.readyRatio * 100)}%`} />
              </div>

              {readiness.blockingReasons.length > 0 && (
                <div className="mt-6">
                  <p className="text-sm uppercase tracking-[0.2em] text-ink/45">阻断原因</p>
                  <ul className="mt-3 list-disc space-y-1 pl-5 text-sm leading-7 text-ink/75">
                    {readiness.blockingReasons.map((reason) => (
                      <li key={reason}>{reason}</li>
                    ))}
                  </ul>
                </div>
              )}
            </>
          )}
        </div>
      )}
    </section>
  );
}
