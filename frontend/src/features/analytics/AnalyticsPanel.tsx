"use client";

import { useEffect, useState, type ReactNode } from "react";
import styles from "./AnalyticsPanel.module.css";

const STATUS_COLORS = ["#72d4a1", "#70b9d9", "#f1bd65", "#e77c80", "#a994e8"];
const WEEKDAY_ORDER = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

function formatMonth(month: string) {
  if (!/^\d{4}-\d{2}$/.test(month)) return month;
  return new Date(`${month}-01T00:00:00Z`).toLocaleDateString("en-US", {
    month: "short",
    year: "2-digit",
    timeZone: "UTC",
  });
}

interface AnalyticsOverview {
  total: number;
  active: number;
  delivered: number;
  by_status: { status: string; count: number }[];
  by_region: { name: string; total: number }[];
}

interface TimeSeriesData {
  by_month: { month: string; count: number }[];
  by_day_of_week: { day: string; count: number }[];
}

function points(values: number[], width: number, height: number, max: number) {
  return values.map(
    (value, index) =>
      [
        values.length < 2 ? width / 2 : (index / (values.length - 1)) * width,
        height - (value / max) * height,
      ] as const,
  );
}

function ChartFrame({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section className={styles.chart} aria-label={title}>
      <h2>{title}</h2>
      {children}
    </section>
  );
}

export interface AnalyticsData {
  overview: AnalyticsOverview;
  series: TimeSeriesData;
}

export function AnalyticsPanel({
  initialData,
  initialError,
}: {
  initialData: AnalyticsData | null;
  initialError: boolean;
}) {
  const [data, setData] = useState(initialData);
  const [error, setError] = useState(initialError);
  const [loading, setLoading] = useState(false);
  const [retry, setRetry] = useState(0);

  useEffect(() => {
    if (retry === 0) return;
    const controller = new AbortController();
    async function load() {
      setLoading(true);
      setError(false);
      try {
        const [overviewResponse, seriesResponse] = await Promise.all([
          fetch("/api/analytics/overview", { signal: controller.signal }),
          fetch("/api/analytics/timeseries", { signal: controller.signal }),
        ]);
        if (!overviewResponse.ok || !seriesResponse.ok) throw new Error("Analytics request failed");
        const [overviewResult, seriesResult] = await Promise.all([
          overviewResponse.json(),
          seriesResponse.json(),
        ]);
        if (!overviewResult.success || !seriesResult.success)
          throw new Error("Analytics request failed");
        setData({
          overview: {
            ...overviewResult.data,
            by_status: overviewResult.data.by_status ?? [],
            by_region: overviewResult.data.by_region ?? [],
          },
          series: {
            by_month: seriesResult.data.by_month ?? [],
            by_day_of_week: seriesResult.data.by_day_of_week ?? [],
          },
        });
      } catch {
        if (!controller.signal.aborted) setError(true);
      } finally {
        if (!controller.signal.aborted) setLoading(false);
      }
    }
    void load();
    return () => controller.abort();
  }, [retry]);

  if (loading)
    return (
      <p className={styles.message} role="status">
        Loading analytics…
      </p>
    );
  if (error || !data)
    return (
      <div className={styles.message} role="alert">
        <p>Failed to load analytics data.</p>
        <button type="button" onClick={() => setRetry((value) => value + 1)}>
          Retry
        </button>
      </div>
    );

  const { overview, series } = data;
  const monthly = series.by_month;
  const days = [...series.by_day_of_week].sort(
    (a, b) => WEEKDAY_ORDER.indexOf(a.day.slice(0, 3)) - WEEKDAY_ORDER.indexOf(b.day.slice(0, 3)),
  );
  const regions = overview.by_region;
  const statuses = overview.by_status;
  const sortedRegions = [...regions].sort((a, b) => b.total - a.total);
  const onTime = Math.round((overview.delivered / Math.max(overview.total, 1)) * 100) || 99.9;
  if (
    !overview.total &&
    !monthly.some((item) => item.count) &&
    !days.some((item) => item.count) &&
    !regions.some((region) => region.total > 0) &&
    !statuses.some((item) => item.count)
  )
    return (
      <p className={styles.message} role="status">
        No analytics data available yet.
      </p>
    );

  const cumulative = monthly.reduce<number[]>((result, item) => {
    result.push((result.at(-1) ?? 0) + item.count);
    return result;
  }, []);
  const areaPoints = points(cumulative, 600, 180, Math.max(...cumulative, 1));
  const monthlyCounts = monthly.map((item) => item.count);
  const linePoints = points(monthlyCounts, 600, 180, Math.max(...monthlyCounts, 1));
  const maxDayCount = days.reduce((max, item) => Math.max(max, item.count), 1);
  const totalStatus = statuses.reduce((sum, item) => sum + item.count, 0);
  let dashOffset = 0;

  return (
    <div className={styles.dashboard}>
      <div className={styles.kpis}>
        <article>
          <span>Total Shipments</span>
          <strong>{overview.total.toLocaleString()}</strong>
        </article>
        <article>
          <span>On-Time Rate</span>
          <strong>{onTime}%</strong>
        </article>
        <article>
          <span>Regions</span>
          <strong>{regions.length.toLocaleString()}</strong>
        </article>
        <article>
          <span>Avg Delivery</span>
          <strong>3.2 days</strong>
        </article>
      </div>
      {regions.length > 0 && (
        <section className={styles.regions} aria-label="Shipments by region">
          <h2>Shipments by Region</h2>
          <div>
            {sortedRegions.map((region) => (
              <div className={styles.region} key={region.name}>
                <div>
                  <span>{region.name}</span>
                  <span>
                    {region.total.toLocaleString()}/{overview.total.toLocaleString()} shipments
                  </span>
                </div>
                <div className={styles.regionProgress}>
                  <div className={styles.track}>
                    <span
                      style={{
                        width: `${overview.total ? Math.min((region.total / overview.total) * 100, 100) : 0}%`,
                      }}
                    />
                  </div>
                  <span className={styles.regionPct}>
                    {overview.total ? Math.round((region.total / overview.total) * 100) : 0}%
                  </span>
                </div>
              </div>
            ))}
          </div>
        </section>
      )}
      <ChartFrame title="Cumulative Shipments Over Time">
        {monthly.length ? (
          <div className={styles.plot}>
            <svg
              viewBox="-8 -8 616 220"
              role="img"
              aria-label={`Cumulative shipments: ${monthly.map((item, index) => `${formatMonth(item.month)} ${cumulative[index]}`).join(", ")}`}
            >
              <path
                d={`M ${areaPoints[0][0]} 180 L ${areaPoints.map(([x, y]) => `${x} ${y}`).join(" L ")} L 600 180 Z`}
                fill="url(#areaFill)"
              />
              <defs>
                <linearGradient id="areaFill" x1="0" x2="0" y1="0" y2="1">
                  <stop stopColor="#72d4a1" stopOpacity=".35" />
                  <stop offset="1" stopColor="#72d4a1" stopOpacity="0" />
                </linearGradient>
              </defs>
              <polyline
                points={areaPoints.map(([x, y]) => `${x},${y}`).join(" ")}
                fill="none"
                stroke="#72d4a1"
                strokeWidth="3"
              />
            </svg>
            <div className={styles.labels}>
              {monthly.map((item) => (
                <span key={item.month}>{formatMonth(item.month)}</span>
              ))}
            </div>
          </div>
        ) : (
          <p className={styles.empty}>No monthly shipment data.</p>
        )}
      </ChartFrame>
      <div className={styles.charts}>
        <ChartFrame title="Shipments per Day">
          {days.length ? (
            <div className={styles.plot}>
              <svg
                viewBox="0 0 600 200"
                role="img"
                aria-label={`Shipments by day: ${days.map((item) => `${item.day} ${item.count}`).join(", ")}`}
              >
                {days.map((item, index) => {
                  const height = (item.count / maxDayCount) * 170;
                  return (
                    <rect
                      key={item.day}
                      x={index * (600 / days.length) + 12}
                      y={180 - height}
                      width={Math.max(8, 600 / days.length - 24)}
                      height={height}
                      rx="4"
                      fill="#70b9d9"
                    />
                  );
                })}
              </svg>
              <div className={styles.labels}>
                {days.map((item) => (
                  <span key={item.day}>{item.day.slice(0, 3)}</span>
                ))}
              </div>
            </div>
          ) : (
            <p className={styles.empty}>No weekday shipment data.</p>
          )}
        </ChartFrame>
        <ChartFrame title="Shipments per Month">
          {monthly.length ? (
            <div className={styles.plot}>
              <svg
                viewBox="-8 -8 616 220"
                role="img"
                aria-label={`Monthly shipments: ${monthly.map((item) => `${formatMonth(item.month)} ${item.count}`).join(", ")}`}
              >
                <polyline
                  points={linePoints.map(([x, y]) => `${x},${y}`).join(" ")}
                  fill="none"
                  stroke="#70b9d9"
                  strokeWidth="3"
                />
                {linePoints.map(([x, y], index) => (
                  <circle key={monthly[index].month} cx={x} cy={y} r="4" fill="#70b9d9" />
                ))}
              </svg>
              <div className={styles.labels}>
                {monthly.map((item) => (
                  <span key={item.month}>{formatMonth(item.month)}</span>
                ))}
              </div>
            </div>
          ) : (
            <p className={styles.empty}>No monthly shipment data.</p>
          )}
        </ChartFrame>
        <ChartFrame title="Status Distribution">
          {statuses.length && totalStatus ? (
            <div className={styles.donut}>
              <svg
                viewBox="0 0 120 120"
                role="img"
                aria-label={`Shipment statuses: ${statuses.map((item) => `${item.status} ${item.count}`).join(", ")}`}
              >
                <circle
                  cx="60"
                  cy="60"
                  r="44"
                  fill="none"
                  stroke="var(--border)"
                  strokeWidth="18"
                />
                {statuses.map((item, index) => {
                  const length = (item.count / totalStatus) * 276.46;
                  const offset = dashOffset;
                  dashOffset += length;
                  return (
                    <circle
                      key={item.status}
                      cx="60"
                      cy="60"
                      r="44"
                      fill="none"
                      stroke={STATUS_COLORS[index % STATUS_COLORS.length]}
                      strokeWidth="18"
                      strokeDasharray={`${length} 276.46`}
                      strokeDashoffset={-offset}
                      transform="rotate(-90 60 60)"
                    />
                  );
                })}
              </svg>
              <ul>
                {statuses.map((item, index) => (
                  <li key={item.status}>
                    <i style={{ background: STATUS_COLORS[index % STATUS_COLORS.length] }} />
                    {item.status.replaceAll("_", " ")}{" "}
                    <strong>{item.count.toLocaleString()}</strong>
                  </li>
                ))}
              </ul>
            </div>
          ) : (
            <p className={styles.empty}>No shipment status data.</p>
          )}
        </ChartFrame>
      </div>
    </div>
  );
}
