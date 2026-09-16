'use client';
import { useEffect, useState } from 'react';
import { apiService } from '@/lib/api';
import { formatDate } from '@/lib/dates';
import { Policy } from '@/types';

const statusClass: Record<string, string> = {
  active: 'badge-green',
  lapsed: 'badge-red',
  cancelled: 'badge-gray',
};

const statusLabel: Record<string, string> = {
  active: 'Aktif',
  lapsed: 'Süresi doldu',
  cancelled: 'İptal',
};

export default function PoliciesPage() {
  const [policies, setPolicies] = useState<Policy[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    apiService.getMyPolicies()
      .then((res) => setPolicies(res.data as Policy[]))
      .catch(() => setError('Poliçeler yüklenirken bir hata oluştu.'))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-[#004C3F]">Poliçelerim</h1>
          <p className="text-[#677471] text-sm mt-1">Aktif ve geçmiş poliçeleriniz</p>
        </div>
        <span className="badge badge-green">{policies.length} Poliçe</span>
      </div>

      <div className="card overflow-hidden">
        {loading ? (
          <div className="flex justify-center py-16"><div className="spinner" /></div>
        ) : error ? (
          <div className="text-center py-16 text-red-500">{error}</div>
        ) : policies.length === 0 ? (
          <div className="text-center py-16 text-[#677471]">
            <span className="text-5xl block mb-4">📋</span>
            <p>Henüz poliçeniz bulunmuyor.</p>
          </div>
        ) : (
          <table className="dash-table">
            <thead>
              <tr>
                <th>Poliçe No</th>
                <th>Tür</th>
                <th>Başlangıç</th>
                <th>Bitiş</th>
                <th>Prim</th>
                <th>Durum</th>
              </tr>
            </thead>
            <tbody>
              {policies.map((p) => (
                <tr key={p.id}>
                  <td className="font-mono text-[#004C3F] font-semibold">{p.policy_number}</td>
                  <td>{p.product_name}</td>
                  <td>{formatDate(p.start_date)}</td>
                  <td>{formatDate(p.end_date)}</td>
                  <td className="font-semibold">{p.premium_amount} {p.currency || '₺'}</td>
                  <td><span className={`badge ${statusClass[p.status] ?? 'badge-gray'}`}>{statusLabel[p.status] ?? p.status}</span></td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
