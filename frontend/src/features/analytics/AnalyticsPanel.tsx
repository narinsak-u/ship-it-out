"use client";

import { useEffect, useState, type ReactNode } from "react";
import {
  Area,
  AreaChart,
  Bar,
  BarChart,
  Cell,
  Line,
  LineChart,
  Pie,
  PieChart,
  XAxis,
} from "recharts";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import { Button } from "@/components/ui/button";
import { getAnalyticsData } from "./client-api";
import type { AnalyticsData } from "./data";

const STATUS_COLORS = [
  "var(--chart-1)",
  "var(--chart-2)",
  "var(--chart-3)",
  "var(--chart-4)",
  "var(--chart-5)",
];
const WEEKDAY_ORDER = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];
const cumulativeConfig = {
  total: { label: "Shipments", color: "var(--chart-1)" },
} satisfies ChartConfig;
const monthlyConfig = {
  count: { label: "Shipments", color: "var(--chart-2)" },
} satisfies ChartConfig;
const weekdayConfig = {
  count: { label: "Shipments", color: "var(--chart-2)" },
} satisfies ChartConfig;
const statusConfig = {
  count: { label: "Shipments", color: "var(--chart-1)" },
} satisfies ChartConfig;

function formatMonth(month: string) {
  if (!/^\d{4}-\d{2}$/.test(month)) return month;
  return new Date(`${month}-01T00:00:00Z`).toLocaleDateString("en-US", {
    month: "short",
    year: "2-digit",
    timeZone: "UTC",
  });
}

function ChartFrame({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section
      className="min-w-0 rounded-[10px] border border-border bg-card p-[18px] shadow-[0_8px_24px_#0002]"
      aria-label={title}
    >
      <h2 className="mb-[18px] font-mono text-[11px] uppercase tracking-[0.1em] text-muted-foreground">
        {title}
      </h2>
      {children}
    </section>
  );
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
        setData(await getAnalyticsData(controller.signal));
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
      <p className="m-0 p-7 text-center font-mono text-[13px] text-muted-foreground" role="status">
        Loading analytics…
      </p>
    );
  if (error || !data)
    return (
      <div className="p-7 text-center font-mono text-[13px] text-muted-foreground" role="alert">
        <p>Failed to load analytics data.</p>
        <Button
          className="mx-auto mt-3 block min-h-11 px-4"
          type="button"
          onClick={() => setRetry((value) => value + 1)}
        >
          Retry
        </Button>
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
      <p className="m-0 p-7 text-center font-mono text-[13px] text-muted-foreground" role="status">
        No analytics data available yet.
      </p>
    );

  const cumulative = monthly.reduce<{ month: string; total: number }[]>((result, item) => {
    result.push({ month: item.month, total: (result.at(-1)?.total ?? 0) + item.count });
    return result;
  }, []);
  const cumulativeLabel = `Cumulative shipments: ${cumulative.map((item) => `${formatMonth(item.month)} ${item.total}`).join(", ")}`;
  const weekdayLabel = `Shipments by day: ${days.map((item) => `${item.day} ${item.count}`).join(", ")}`;
  const monthlyLabel = `Monthly shipments: ${monthly.map((item) => `${formatMonth(item.month)} ${item.count}`).join(", ")}`;
  const statusLabel = `Shipment statuses: ${statuses.map((item) => `${item.status} ${item.count}`).join(", ")}`;
  const labelsClass =
    "flex justify-between gap-1 overflow-hidden font-mono text-[10px] text-muted-foreground [&>span]:min-w-0 [&>span]:overflow-hidden [&>span]:text-ellipsis [&>span]:whitespace-nowrap";
  const emptyClass = "m-0 p-7 text-center font-mono text-[13px] text-muted-foreground";

  return (
    <div className="grid gap-6">
      <div className="grid grid-cols-2 gap-[10px] min-[621px]:grid-cols-4 min-[621px]:gap-[14px]">
        <article className="grid min-w-0 gap-3 rounded-[10px] border border-border bg-card p-[18px] shadow-[0_8px_24px_#0002]">
          <span className="font-mono text-[11px] uppercase tracking-[0.1em] text-muted-foreground">
            Total Shipments
          </span>
          <strong className="font-mono text-[clamp(22px,3vw,30px)] font-semibold">
            {overview.total.toLocaleString()}
          </strong>
        </article>
        <article className="grid min-w-0 gap-3 rounded-[10px] border border-border bg-card p-[18px] shadow-[0_8px_24px_#0002]">
          <span className="font-mono text-[11px] uppercase tracking-[0.1em] text-muted-foreground">
            On-Time Rate
          </span>
          <strong className="font-mono text-[clamp(22px,3vw,30px)] font-semibold text-success">
            {onTime}%
          </strong>
        </article>
        <article className="grid min-w-0 gap-3 rounded-[10px] border border-border bg-card p-[18px] shadow-[0_8px_24px_#0002]">
          <span className="font-mono text-[11px] uppercase tracking-[0.1em] text-muted-foreground">
            Regions
          </span>
          <strong className="font-mono text-[clamp(22px,3vw,30px)] font-semibold text-info">
            {regions.length.toLocaleString()}
          </strong>
        </article>
        <article className="grid min-w-0 gap-3 rounded-[10px] border border-border bg-card p-[18px] shadow-[0_8px_24px_#0002]">
          <span className="font-mono text-[11px] uppercase tracking-[0.1em] text-muted-foreground">
            Avg Delivery
          </span>
          <strong className="font-mono text-[clamp(22px,3vw,30px)] font-semibold">3.2 days</strong>
        </article>
      </div>
      {regions.length > 0 && (
        <section className="grid gap-3" aria-label="Shipments by region">
          <h2 className="m-0 font-mono text-sm font-semibold">Shipments by Region</h2>
          <div className="grid gap-3">
            {sortedRegions.map((region) => (
              <div
                className="grid gap-[10px] rounded-lg border border-border bg-card p-[14px] font-mono text-xs"
                key={region.name}
              >
                <div className="flex justify-between gap-2">
                  <span>{region.name}</span>
                  <span className="text-muted-foreground">
                    {region.total.toLocaleString()}/{overview.total.toLocaleString()} shipments
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  <div className="h-1.5 flex-1 overflow-hidden rounded-md bg-secondary">
                    <span
                      className="block h-full rounded-md bg-gradient-to-br from-accent to-primary"
                      style={{
                        width: `${overview.total ? Math.min((region.total / overview.total) * 100, 100) : 0}%`,
                      }}
                    />
                  </div>
                  <span className="whitespace-nowrap font-mono text-[11px] text-muted-foreground">
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
          <>
            <ChartContainer
              config={cumulativeConfig}
              role="img"
              aria-label={cumulativeLabel}
              className="h-[190px] w-full aspect-auto"
            >
              <AreaChart data={cumulative} margin={{ top: 8, right: 8, bottom: 0, left: 8 }}>
                <XAxis dataKey="month" hide />
                <ChartTooltip
                  content={
                    <ChartTooltipContent
                      labelFormatter={(_, payload) =>
                        formatMonth(String(payload?.[0]?.payload?.month ?? ""))
                      }
                    />
                  }
                />
                <Area
                  dataKey="total"
                  type="monotone"
                  stroke="#72d4a1"
                  strokeWidth={3}
                  fill="#72d4a1"
                  fillOpacity={0.2}
                />
              </AreaChart>
            </ChartContainer>
            <div className={labelsClass}>
              {monthly.map((item) => (
                <span key={item.month}>{formatMonth(item.month)}</span>
              ))}
            </div>
          </>
        ) : (
          <p className={emptyClass}>No monthly shipment data.</p>
        )}
      </ChartFrame>
      <div className="grid grid-cols-1 gap-[14px] min-[621px]:grid-cols-2 min-[901px]:grid-cols-3 [&>:last-child]:min-[901px]:col-span-1 [&>:last-child]:min-[621px]:col-span-full">
        <ChartFrame title="Shipments per Day">
          {days.length ? (
            <>
              <ChartContainer
                config={weekdayConfig}
                role="img"
                aria-label={weekdayLabel}
                className="h-[190px] w-full aspect-auto max-[620px]:h-40"
              >
                <BarChart data={days} margin={{ top: 8, right: 8, bottom: 0, left: 8 }}>
                  <XAxis dataKey="day" hide />
                  <ChartTooltip
                    content={<ChartTooltipContent labelFormatter={(label) => String(label)} />}
                  />
                  <Bar dataKey="count" fill="#70b9d9" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ChartContainer>
              <div className={labelsClass}>
                {days.map((item) => (
                  <span key={item.day}>{item.day.slice(0, 3)}</span>
                ))}
              </div>
            </>
          ) : (
            <p className={emptyClass}>No weekday shipment data.</p>
          )}
        </ChartFrame>
        <ChartFrame title="Shipments per Month">
          {monthly.length ? (
            <>
              <ChartContainer
                config={monthlyConfig}
                role="img"
                aria-label={monthlyLabel}
                className="h-[190px] w-full aspect-auto max-[620px]:h-40"
              >
                <LineChart data={monthly} margin={{ top: 8, right: 8, bottom: 0, left: 8 }}>
                  <XAxis dataKey="month" hide />
                  <ChartTooltip
                    content={
                      <ChartTooltipContent labelFormatter={(label) => formatMonth(String(label))} />
                    }
                  />
                  <Line
                    dataKey="count"
                    type="monotone"
                    stroke="#70b9d9"
                    strokeWidth={3}
                    dot={{ r: 4, fill: "#70b9d9" }}
                  />
                </LineChart>
              </ChartContainer>
              <div className={labelsClass}>
                {monthly.map((item) => (
                  <span key={item.month}>{formatMonth(item.month)}</span>
                ))}
              </div>
            </>
          ) : (
            <p className={emptyClass}>No monthly shipment data.</p>
          )}
        </ChartFrame>
        <ChartFrame title="Status Distribution">
          {statuses.length && statuses.some((item) => item.count) ? (
            <div className="grid grid-cols-[minmax(90px,140px)_1fr] items-center gap-3 max-[620px]:grid-cols-[minmax(88px,120px)_1fr]">
              <ChartContainer
                config={statusConfig}
                role="img"
                aria-label={statusLabel}
                className="w-full aspect-square"
              >
                <PieChart>
                  <ChartTooltip content={<ChartTooltipContent labelKey="status" />} />
                  <Pie
                    data={statuses}
                    dataKey="count"
                    nameKey="status"
                    innerRadius="55%"
                    outerRadius="85%"
                    strokeWidth={0}
                  >
                    {statuses.map((item, index) => (
                      <Cell key={item.status} fill={STATUS_COLORS[index % STATUS_COLORS.length]} />
                    ))}
                  </Pie>
                </PieChart>
              </ChartContainer>
              <ul className="m-0 grid list-none gap-2 p-0 font-mono text-[11px] capitalize">
                {statuses.map((item, index) => (
                  <li className="flex items-center gap-1.5" key={item.status}>
                    <i
                      className="h-2 w-2 shrink-0 rounded-full"
                      style={{ background: STATUS_COLORS[index % STATUS_COLORS.length] }}
                    />
                    {item.status.replaceAll("_", " ")}{" "}
                    <strong className="ml-auto">{item.count.toLocaleString()}</strong>
                  </li>
                ))}
              </ul>
            </div>
          ) : (
            <p className={emptyClass}>No shipment status data.</p>
          )}
        </ChartFrame>
      </div>
    </div>
  );
}
