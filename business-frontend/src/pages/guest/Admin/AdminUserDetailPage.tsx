import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { apiClient } from '@/api/client';
import { getTokenRole } from '@/utils/auth';
import type { FullUserDetails } from '@/types/api';
import styles from './AdminUserDetailPage.module.css';

const STATUS_OPTIONS = ['active', 'muted', 'suspended'] as const;
const ROLE_OPTIONS = ['user', 'business', 'moderator', 'admin'] as const;

export const AdminUserDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<FullUserDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const [statusSaving, setStatusSaving] = useState(false);
  const [statusMsg, setStatusMsg] = useState('');
  const [statusError, setStatusError] = useState('');

  const [roleSaving, setRoleSaving] = useState(false);
  const [roleMsg, setRoleMsg] = useState('');
  const [roleError, setRoleError] = useState('');
  const [selectedRole, setSelectedRole] = useState('');

  const tokenRole = getTokenRole();
  const isAdmin = tokenRole === 'admin';

  useEffect(() => {
    if (!id) return;
    const load = async () => {
      try {
        const data = await apiClient.adminGetUserDetails(id);
        setUser(data);
        setSelectedRole(data.role_slug);
      } catch (err: any) {
        setError(err?.response?.data?.error || err?.message || 'Failed to load user.');
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [id]);

  const handleStatusChange = async (newStatus: string) => {
    if (!user || !id) return;
    setStatusSaving(true);
    setStatusMsg('');
    setStatusError('');
    try {
      await apiClient.updateUserStatus(id, newStatus);
      setUser(prev => prev ? { ...prev, status: newStatus } : prev);
      setStatusMsg(`Status updated to "${newStatus}".`);
    } catch (err: any) {
      setStatusError(err?.response?.data?.error || err?.message || 'Failed to update status.');
    } finally {
      setStatusSaving(false);
    }
  };

  const handleRoleChange = async () => {
    if (!user || !id || selectedRole === user.role_slug) return;
    setRoleSaving(true);
    setRoleMsg('');
    setRoleError('');
    try {
      await apiClient.adminChangeUserRole(id, selectedRole);
      setUser(prev => prev ? { ...prev, role_slug: selectedRole } : prev);
      setRoleMsg(`Role updated to "${selectedRole}".`);
    } catch (err: any) {
      setRoleError(err?.response?.data?.error || err?.message || 'Failed to update role.');
    } finally {
      setRoleSaving(false);
    }
  };

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

  if (loading) {
    return (
      <div className={styles.page}>
        <div className={styles.container}>
          <div className={styles.loadingRow}><div className={styles.spinner} /></div>
        </div>
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className={styles.page}>
        <div className={styles.container}>
          <Link to="/admin" className={styles.back}>← Back to users</Link>
          <p className={styles.errorBox}>{error || 'User not found.'}</p>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <div className={styles.container}>
        <Link to="/admin" className={styles.back}>← Back to users</Link>

        {/* Identity card */}
        <section className={styles.card}>
          <h1 className={styles.email}>{user.email}</h1>
          <div className={styles.badges}>
            <span className={`${styles.badge} ${roleBadgeClass(user.role_slug)}`}>{user.role_slug}</span>
            <span className={`${styles.badge} ${statusBadgeClass(user.status)}`}>{user.status}</span>
          </div>
          <p className={styles.meta}>
            ID: <code>{user.id}</code>
          </p>
          <p className={styles.meta}>
            Joined: {new Date(user.created_at).toLocaleString()}
          </p>
        </section>

        {/* Customer profile */}
        {user.customer_profile && (
          <section className={styles.card}>
            <h2 className={styles.sectionTitle}>Customer Profile</h2>
            <div className={styles.profileGrid}>
              <div className={styles.profileField}>
                <span className={styles.fieldLabel}>Username</span>
                <span>{user.customer_profile.username}</span>
              </div>
              <div className={styles.profileField}>
                <span className={styles.fieldLabel}>Name</span>
                <span>{user.customer_profile.first_name} {user.customer_profile.last_name}</span>
              </div>
              {user.customer_profile.bio && (
                <div className={styles.profileField}>
                  <span className={styles.fieldLabel}>Bio</span>
                  <span>{user.customer_profile.bio}</span>
                </div>
              )}
            </div>
          </section>
        )}

        {/* Business profile */}
        {user.business_profile && (
          <section className={styles.card}>
            <h2 className={styles.sectionTitle}>Business Profile</h2>
            <div className={styles.profileGrid}>
              <div className={styles.profileField}>
                <span className={styles.fieldLabel}>Name</span>
                <span>{user.business_profile.name}</span>
              </div>
              <div className={styles.profileField}>
                <span className={styles.fieldLabel}>Type</span>
                <span className={styles.capitalize}>{user.business_profile.profile_type}</span>
              </div>
              {user.business_profile.description && (
                <div className={styles.profileField}>
                  <span className={styles.fieldLabel}>Description</span>
                  <span>{user.business_profile.description}</span>
                </div>
              )}
              {user.business_profile.latitude != null && (
                <div className={styles.profileField}>
                  <span className={styles.fieldLabel}>Location</span>
                  <span>{user.business_profile.latitude}, {user.business_profile.longitude}</span>
                </div>
              )}
            </div>
          </section>
        )}

        {/* Status change */}
        <section className={styles.card}>
          <h2 className={styles.sectionTitle}>Account Status</h2>
          <p className={styles.sectionHint}>
            Current: <span className={`${styles.badge} ${statusBadgeClass(user.status)}`}>{user.status}</span>
          </p>
          <div className={styles.actionRow}>
            {STATUS_OPTIONS.map(opt => (
              <button
                key={opt}
                className={`${styles.statusBtn} ${user.status === opt ? styles.statusBtnActive : ''}`}
                disabled={statusSaving || user.status === opt}
                onClick={() => handleStatusChange(opt)}
              >
                {opt}
              </button>
            ))}
          </div>
          {statusMsg && <p className={styles.successMsg}>{statusMsg}</p>}
          {statusError && <p className={styles.errorMsg}>{statusError}</p>}
        </section>

        {/* Role change — admin only */}
        {isAdmin && (
          <section className={styles.card}>
            <h2 className={styles.sectionTitle}>Change Role</h2>
            <p className={styles.sectionHint}>
              Current: <span className={`${styles.badge} ${roleBadgeClass(user.role_slug)}`}>{user.role_slug}</span>
            </p>
            <div className={styles.roleRow}>
              <select
                className={styles.select}
                value={selectedRole}
                onChange={e => setSelectedRole(e.target.value)}
                disabled={roleSaving}
              >
                {ROLE_OPTIONS.map(r => (
                  <option key={r} value={r}>{r}</option>
                ))}
              </select>
              <button
                className={styles.saveRoleBtn}
                onClick={handleRoleChange}
                disabled={roleSaving || selectedRole === user.role_slug}
              >
                {roleSaving ? 'Saving…' : 'Save Role'}
              </button>
            </div>
            {roleMsg && <p className={styles.successMsg}>{roleMsg}</p>}
            {roleError && <p className={styles.errorMsg}>{roleError}</p>}
          </section>
        )}
      </div>
    </div>
  );
};
