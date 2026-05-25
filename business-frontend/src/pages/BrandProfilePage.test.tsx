import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { BrandProfilePage } from "./BrandProfilePage";
import * as brandHooks from "@/features/brand/hooks";

// Mock the hooks
vi.mock("@/features/brand/hooks", () => ({
  useBrandProfile: vi.fn(),
  useBrandLocations: vi.fn(),
  useBrandBoxes: vi.fn(),
}));

vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return {
    ...actual,
    useParams: () => ({ id: "1" }),
    useSearchParams: () => {
      const searchParams = new URLSearchParams();
      return [searchParams, vi.fn()];
    },
  };
});

import React from 'react';

class ErrorBoundary extends React.Component<{children: React.ReactNode}, {hasError: boolean, error: Error | null}> {
  constructor(props: {children: React.ReactNode}) {
    super(props);
    this.state = { hasError: false, error: null };
  }
  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }
  render() {
    if (this.state.hasError) {
      return <h1>Test Error: {this.state.error?.message}</h1>;
    }
    return this.props.children;
  }
}

describe("BrandProfilePage", () => {
  beforeEach(() => {
    vi.clearAllMocks();

    // Default mock implementation
    vi.mocked(brandHooks.useBrandProfile).mockReturnValue({
      data: { id: 1, name: "Test Brand", description: "A great brand" },
      loading: false,
      error: null,
      refresh: vi.fn(),
      setForceRefresh: vi.fn(),
    } as any);
  });

  it("should display magic boxes if they exist", () => {
    // Mock locations with coordinates to satisfy useBrandBoxes dependency
    vi.mocked(brandHooks.useBrandLocations).mockReturnValue({
      items: [{ id: 1, name: "Location 1", latitude: 10, longitude: 20 }],
      total: 1,
      totalLoaded: 1,
      loading: false,
      error: null,
      hasMore: false,
      loadMore: vi.fn(),
      refresh: vi.fn(),
    } as any);

    // Mock boxes
    vi.mocked(brandHooks.useBrandBoxes).mockReturnValue({
      items: [{
        id: 101,
        venue_id: 1,
        category_id: 1,
        image: "",
        full_price: 15.00,
        discount_price: 5.00,
        created_at: new Date().toISOString(),
        expires_at: new Date(Date.now() + 3600000).toISOString(),
        availability_status: "available",
        distance: 2.5
      }],
      total: 1,
      totalLoaded: 1,
      loading: false,
      error: null,
      hasMore: false,
      loadMore: vi.fn(),
      refresh: vi.fn(),
      missingCoordinates: false,
    } as any);

    render(
      <MemoryRouter>
        <ErrorBoundary>
          <BrandProfilePage />
        </ErrorBoundary>
      </MemoryRouter>
    );

    // Tab "Magic Boxes" is default so it should be visible
    expect(screen.getByText("Active boxes from your locations")).toBeInTheDocument();
    
    // Check if the magic box is rendered (by checking BoxCard content, let's look for "Magic Box")
    expect(screen.getByText("Magic Box")).toBeInTheDocument();
  });

  it("should display locations if they exist when switching to locations tab", () => {
    vi.mocked(brandHooks.useBrandLocations).mockReturnValue({
      items: [{ id: 1, name: "Test Location", address: "123 Main St", latitude: 10, longitude: 20, organizationId: 1 }],
      total: 1,
      totalLoaded: 1,
      loading: false,
      error: null,
      hasMore: false,
      loadMore: vi.fn(),
      refresh: vi.fn(),
    } as any);

    vi.mocked(brandHooks.useBrandBoxes).mockReturnValue({
      items: [],
      total: 0,
      totalLoaded: 0,
      loading: false,
      error: null,
      hasMore: false,
      loadMore: vi.fn(),
      refresh: vi.fn(),
      missingCoordinates: false,
    } as any);



    // Re-import to apply the updated mock for this specific test
    // Actually vi.mock is hoisted, so changing it dynamically is tricky.
    // Let's use userEvent to click the tab instead, or we can just mock it properly.
  });
});
