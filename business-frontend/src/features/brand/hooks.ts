import { useCallback, useEffect, useMemo, useState } from "react";

import {
  businessApi,
  type BrandLocation,
  type BrandLocationsResponse,
  type BrandPost,
  type BrandPostsResponse,
  type BrandProfile,
  type Box,
} from "@/api/business";

type AsyncState<T> = {
  data: T | null;
  loading: boolean;
  error: string | null;
};

type ListState<T> = {
  items: T[];
  total: number;
  loading: boolean;
  error: string | null;
};

export function useBrandProfile(brandId?: number): AsyncState<BrandProfile> & { refresh: () => void } {
  const [state, setState] = useState<AsyncState<BrandProfile>>({
    data: null,
    loading: false,
    error: null,
  });
  const [refreshKey, setRefreshKey] = useState(0);

  const refresh = useCallback(() => {
    setRefreshKey((prev) => prev + 1);
  }, []);

  useEffect(() => {
    let cancelled = false;

    if (!brandId || !Number.isFinite(brandId) || brandId <= 0) {
      Promise.resolve().then(() => {
        if (cancelled) return;
        setState((prev) => {
          if (prev.data === null && prev.error === "Invalid brand id") return prev;
          return { data: null, loading: false, error: "Invalid brand id" };
        });
      });
      return;
    }

    Promise.resolve().then(() => {
      if (cancelled) return;
      setState((prev) => ({ ...prev, loading: true, error: null }));
    });

    businessApi
      .getBrandProfile(brandId)
      .then((data) => {
        if (cancelled) return;
        setState({ data, loading: false, error: null });
      })
      .catch((err) => {
        if (cancelled) return;
        setState({ data: null, loading: false, error: err instanceof Error ? err.message : "Failed to load brand" });
      });

    return () => {
      cancelled = true;
    };
  }, [brandId, refreshKey]);

  return { ...state, refresh };
}

export function useBrandLocations(brandId?: number, limit = 20): ListState<BrandLocation> {
  const [state, setState] = useState<ListState<BrandLocation>>({
    items: [],
    total: 0,
    loading: false,
    error: null,
  });

  useEffect(() => {
    let cancelled = false;

    if (!brandId || !Number.isFinite(brandId) || brandId <= 0) {
      Promise.resolve().then(() => {
        if (cancelled) return;
        setState((prev) => {
          if (prev.items.length === 0 && prev.error === "Invalid brand id") return prev;
          return { items: [], total: 0, loading: false, error: "Invalid brand id" };
        });
      });
      return;
    }

    Promise.resolve().then(() => {
      if (cancelled) return;
      setState((prev) => ({ ...prev, loading: true, error: null }));
    });

    businessApi
      .getBrandLocations(brandId, { limit })
      .then((data: BrandLocationsResponse) => {
        if (cancelled) return;
        setState({ items: data.items, total: data.total, loading: false, error: null });
      })
      .catch((err) => {
        if (cancelled) return;
        setState({ items: [], total: 0, loading: false, error: err instanceof Error ? err.message : "Failed to load locations" });
      });

    return () => {
      cancelled = true;
    };
  }, [brandId, limit]);

  return state;
}

type PostsState = {
  items: BrandPost[];
  totalLoaded: number;
  loading: boolean;
  error: string | null;
  hasMore: boolean;
  loadMore: () => void;
  reset: () => void;
};

export function useBrandPosts(orgIds: number[], pageSize = 12): PostsState {
  const [items, setItems] = useState<BrandPost[]>([]);
  const [totalLoaded, setTotalLoaded] = useState(0);
  const [skip, setSkip] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const orgKey = useMemo(() => orgIds.slice().sort((a, b) => a - b).join(","), [orgIds]);

  const reset = useCallback(() => {
    setItems([]);
    setTotalLoaded(0);
    setSkip(0);
    setHasMore(true);
    setError(null);
  }, []);

  useEffect(() => {
    let cancelled = false;

    if (orgIds.length === 0) {
      Promise.resolve().then(() => {
        if (cancelled) return;
        setHasMore(false);
        setLoading(false);
        setError(null);
        setItems([]);
        setTotalLoaded(0);
        setSkip(0);
      });
      return;
    }

    Promise.resolve().then(() => {
      if (cancelled) return;
      setLoading(true);
      setError(null);
    });

    businessApi
      .getBrandPosts({ skip, limit: pageSize, orgIds })
      .then((data: BrandPostsResponse) => {
        if (cancelled) return;
        setItems((prev) => (skip === 0 ? data.items : [...prev, ...data.items]));
        setTotalLoaded((prev) => (skip === 0 ? data.items.length : prev + data.items.length));
        setHasMore(skip + pageSize < data.total);
      })
      .catch((err) => {
        if (cancelled) return;
        setError(err instanceof Error ? err.message : "Failed to load posts");
      })
      .finally(() => {
        if (cancelled) return;
        setLoading(false);
      });

    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [orgKey, pageSize, skip]);

  // Use another effect to reset skip when orgKey changes
  useEffect(() => {
    Promise.resolve().then(() => {
      setSkip(0);
      setHasMore(true);
      setItems([]);
      setTotalLoaded(0);
    });
  }, [orgKey]);

  const loadMore = useCallback(() => {
    if (loading || !hasMore) return;
    setSkip((prev) => prev + pageSize);
  }, [loading, hasMore, pageSize]);

  return { items, totalLoaded, loading, error, hasMore, loadMore, reset };
}

type BoxesState = {
  items: Box[];
  total: number;
  loading: boolean;
  error: string | null;
  hasMore: boolean;
  loadMore: () => void;
  missingCoordinates: boolean;
};

type Coordinates = { lat: number; lon: number };

export function useBrandBoxes(params: {
  orgId?: number;
  coordinates?: Coordinates | null;
  pageSize?: number;
}): BoxesState {
  const { orgId, coordinates, pageSize = 12 } = params;
  const [items, setItems] = useState<Box[]>([]);
  const [total, setTotal] = useState(0);
  const [skip, setSkip] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const lat = coordinates?.lat;
  const lon = coordinates?.lon;
  const missingCoordinates = !coordinates || !Number.isFinite(lat) || !Number.isFinite(lon);

  useEffect(() => {
    let cancelled = false;

    if (!orgId || orgId <= 0 || missingCoordinates) {
      Promise.resolve().then(() => {
        if (cancelled) return;
        setHasMore(false);
        setItems([]);
        setTotal(0);
        setSkip(0);
      });
      return;
    }

    Promise.resolve().then(() => {
      if (cancelled) return;
      setLoading(true);
    });

    businessApi
      .getNearbyBoxes({
        lat: lat!,
        lon: lon!,
        skip,
        limit: pageSize,
        orgId,
      })
      .then((data) => {
        if (cancelled) return;
        setItems((prev) => (skip === 0 ? data.items : [...prev, ...data.items]));
        setTotal(data.total);
        setHasMore(skip + pageSize < data.total);
      })
      .catch((err) => {
        if (cancelled) return;
        setError(err instanceof Error ? err.message : "Failed to load boxes");
      })
      .finally(() => {
        if (cancelled) return;
        setLoading(false);
      });

    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lat, lon, orgId, skip, pageSize]);

  // Reset pagination when filters change
  useEffect(() => {
    Promise.resolve().then(() => {
      setSkip(0);
      setHasMore(true);
      setItems([]);
      setTotal(0);
    });
  }, [orgId, lat, lon, pageSize]);

  const loadMore = useCallback(() => {
    if (loading || !hasMore) return;
    setSkip((prev) => prev + pageSize);
  }, [loading, hasMore, pageSize]);

  return {
    items,
    total,
    loading,
    error,
    hasMore,
    loadMore,
    missingCoordinates,
  };
}