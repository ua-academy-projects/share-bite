import { useEffect, useRef, useState } from 'react';

export interface LabNotification {
    id: string;
    type: string;
    entityID: string;
    metadata?: any;
    isRead: boolean;
    createdAt: string;
    readAt?: string;
}

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected';

export function useRealtimeNotifications(token: string | null) {
    const [notifications, setNotifications] = useState<LabNotification[]>([]);
    const [status, setStatus] = useState<ConnectionStatus>('disconnected');
    const eventSourceRef = useRef<EventSource | null>(null);
    const reconnectTimeoutRef = useRef<number | null>(null);

    useEffect(() => {
        if (!token) {
            setStatus('disconnected');
            setNotifications([]);
            return;
        }

        const connect = () => {
            if (eventSourceRef.current) {
                eventSourceRef.current.close();
            }

            setStatus('connecting');
            const url = `/api/notifications/stream?access_token=${encodeURIComponent(token)}`;
            const es = new EventSource(url);

            es.onopen = () => {
                console.log("[SSE] Connection established");
                setStatus('connected');
            };

            es.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data) as LabNotification;
                    setNotifications((prev) => [data, ...prev].slice(0, 100));
                } catch (err) {
                    console.error("[SSE] Failed to parse event data", err);
                }
            };

            es.onerror = () => {
                console.error("[SSE] Connection error, retrying in 5s");
                setStatus('connecting');
                es.close();

                if (reconnectTimeoutRef.current) {
                    clearTimeout(reconnectTimeoutRef.current);
                }
                reconnectTimeoutRef.current = window.setTimeout(connect, 5000);
            };

            eventSourceRef.current = es;
        };

        connect();

        return () => {
            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current);
            }
            if (eventSourceRef.current) {
                eventSourceRef.current.close();
            }
        };
    }, [token]);

    const clearNotifications = () => setNotifications([]);

    return {
        notifications,
        setNotifications,
        status,
        clearNotifications
    };
}
