export const orderKeys = {
  all: ["orders"] as const,
  lists: () => [...orderKeys.all, "list"] as const,
  list: (filters: Record<string, unknown>) => [...orderKeys.lists(), filters] as const,
  details: () => [...orderKeys.all, "detail"] as const,
  detail: (id: string) => [...orderKeys.details(), id] as const,
};

export const deliveryKeys = {
  all: ["deliveries"] as const,
  active: () => [...deliveryKeys.all, "active"] as const,
};

export const hubKeys = {
  all: ["hubs"] as const,
};

export const analyticsKeys = {
  all: ["analytics"] as const,
  timeseries: () => [...analyticsKeys.all, "timeseries"] as const,
};

export const eventKeys = {
  all: ["events"] as const,
  byTracking: (tn: string) => [...eventKeys.all, tn] as const,
};
