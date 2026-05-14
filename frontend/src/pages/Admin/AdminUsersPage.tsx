import React, { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { apiClient } from '../../api/client';
import type { AdminUserListItem } from '../../types/api';
import styles from './AdminUsersPage.module.css';

const PAGE_SIZE = 10;

const ROLE_OPTIONS = ['', 'user', 'business', 'moderator', 'admin'];
const STATUS_OPTIONS = ['', 'active', 'muted', 'suspended'];

export const AdminUsersPage: React.FC = () => {
  const [users, setUsers] = useState<AdminUserListItem[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const [search, setSearch] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');
  const [roleFilter, setRoleFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [page, setPage] = useState(0);
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  useEffect(() => {
    const t = setTimeout(() => setDebouncedSearch(search), 350);
    return () => clearTimeout(t);
  }, [search]);
  useEffect(() => {
    setPage(0);
  }, [debouncedSearch, roleFilter, statusFilter, sortOrder]);
  const toggleSort = () => {
    setSortOrder(prev => (prev === 'asc' ? 'desc' : 'asc'));
  };

  const load = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const data = await apiClient.adminGetUsers({
        limit: PAGE_SIZE,
        offset: page * PAGE_SIZE,
        search_email: debouncedSearch || undefined,
        role: roleFilter || undefined,
        status: statusFilter || undefined,
        sort_order: sortOrder,
      });
      setUsers(data.items ?? []);
      setTotalCount(data.total_count ?? 0);
    } catch (err: any) {
      setError(err?.response?.data?.error || err?.message || 'Failed to load users.');
    } finally {
      setLoading(false);
    }
  }, [page, debouncedSearch, roleFilter, statusFilter, sortOrder]);

  useEffect(() => { load(); }, [load]);

  const totalPages = Math.ceil(totalCount / PAGE_SIZE);

  const roleBadgeClass = (role: string) => {
    switch (role) {
      case 'admin': return styles.badgeAdmin;
      case 'moderator': return styles.badgeModerator;
      case 'business': return styles.badgeBusiness;
      default: return styles.badgeUser;
    }
  };

  const statusBadgeClass = (status: string) => {
    switch (status) {
      case 'active': return styles.statusActive;
      case 'muted': return styles.statusMuted;
      case 'suspended': return styles.statusSuspended;
      default: return styles.statusDefault;
    }
  };

  return (
      <div className={styles.page}>
        <div className={styles.container}>
          <div className={styles.header}>
            <div>
              <h1 className={styles.title}>Admin — Users</h1>
              <p className={styles.subtitle}>{totalCount} total users</p>
            </div>
          </div>

          <div className={styles.filters}>
            <input
                className={styles.searchInput}
                type="text"
                placeholder="Search by email…"
                value={search}
                onChange={e => setSearch(e.target.value)}
            />
            <select
                className={styles.select}
                value={roleFilter}
                onChange={e => setRoleFilter(e.target.value)}
            >
              <option value="">All roles</option>
              {ROLE_OPTIONS.filter(Boolean).map(r => (
                  <option key={r} value={r}>{r}</option>
              ))}
            </select>
            <select
                className={styles.select}
                value={statusFilter}
                onChange={e => setStatusFilter(e.target.value)}
            >
              <option value="">All statuses</option>
              {STATUS_OPTIONS.filter(Boolean).map(s => (
                  <option key={s} value={s}>{s}</option>
              ))}
            </select>
          </div>

          {error && <p className={styles.error}>{error}</p>}

          {loading ? (
              <div className={styles.loadingRow}>
                <div className={styles.spinner} />
              </div>
          ) : users.length === 0 ? (
              <div className={styles.empty}>No users found.</div>
          ) : (
              <div className={styles.tableWrapper}>
                <table className={styles.table}>
                  <colgroup>
                    <col className={styles.colEmail} />
                    <col className={styles.colRole} />
                    <col className={styles.colStatus} />
                    <col className={styles.colCreated} />
                    <col className={styles.colActions} />
                  </colgroup>
                  <thead>
                  <tr>
                    <th className={styles.thLeft}>Email</th>
                    <th className={styles.thLeft}>Role</th>
                    <th className={styles.thLeft}>Status</th>
                    <th
                        className={`${styles.thLeft} ${styles.sortableHeader}`}
                        onClick={toggleSort}
                    >
                      Created <span className={styles.sortIcon}>{sortOrder === 'asc' ? '↑' : '↓'}</span>
                    </th>
                    <th className={styles.thRight} scope="col">
                      <span className={styles.visuallyHidden}>Actions</span>
                    </th>
                  </tr>
                  </thead>
                  <tbody>
                  {users.map(user => (
                      <tr key={user.id}>
                        <td className={`${styles.tdLeft} ${styles.emailCell}`}>{user.email}</td>
                        <td className={styles.tdLeft}>
                      <span className={`${styles.badge} ${roleBadgeClass(user.role_slug)}`}>
                        {user.role_slug}
                      </span>
                        </td>
                        <td className={styles.tdLeft}>
                      <span className={`${styles.badge} ${statusBadgeClass(user.status)}`}>
                        {user.status}
                      </span>
                        </td>
                        <td className={`${styles.tdLeft} ${styles.dateCell}`}>
                          {new Date(user.created_at).toLocaleDateString()}
                        </td>
                        <td className={styles.tdRight}>
                          <Link to={`/admin/users/${user.id}`} className={styles.detailsLink}>
                            Details →
                          </Link>
                        </td>
                      </tr>
                  ))}
                  </tbody>
                </table>
              </div>
          )}

          {totalPages > 1 && (
              <div className={styles.pagination}>
                <button
                    className={styles.pageBtn}
                    onClick={() => setPage(p => Math.max(0, p - 1))}
                    disabled={page === 0}
                >
                  ← Prev
                </button>
                <span className={styles.pageInfo}>
              Page {page + 1} of {totalPages}
            </span>
                <button
                    className={styles.pageBtn}
                    onClick={() => setPage(p => Math.min(totalPages - 1, p + 1))}
                    disabled={page >= totalPages - 1}
                >
                  Next →
                </button>
              </div>
          )}
        </div>
      </div>
  );
};