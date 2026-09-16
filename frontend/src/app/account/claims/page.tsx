'use client';
import { useEffect, useState } from 'react';
import { apiService } from '@/lib/api';
import { formatDate } from '@/lib/dates';
import { Claim } from '@/types';

const statusClass: Record<string, string> = {
  open: 'badge-yellow',
  review: 'badge-yellow',
  closed: 'badge-gray',
};

const statusLabel: Record<string, string> = {
  open: 'Açık',
  review: 'İncelemede',
  closed: 'Kapandı',
};

export default function ClaimsPage() {
  const [claims, setClaims] = useState<Claim[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    apiService.getMyClaims()
      .then((res) => setClaims(res.data as Claim[]))
      .catch(() => setError('Hasar kayıtları yüklenirken bir hata oluştu.'))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-[#004C3F]">Hasar Kayıtlarım</h1>
          <p className="text-[#677471] text-sm mt-1">Bildirilen hasar ve talep kayıtlarınız</p>
        </div>
        <span className="badge badge-green">{claims.length} Kayıt</span>
      </div>

      <div className="card overflow-hidden">
        {loading ? (
          <div className="flex justify-center py-16"><div className="spinner" /></div>
        ) : error ? (
          <div className="text-center py-16 text-red-500">{error}</div>
        ) : claims.length === 0 ? (
          <div className="text-center py-16 text-[#677471]">
            <span className="text-5xl block mb-4">🔧</span>
            <p>Hasar kaydınız bulunmuyor.</p>
          </div>
        ) : (
          <table className="dash-table">
            <thead>
              <tr>
                <th>No</th>
                <th>Hasar No</th>
                <th>Açıklama</th>
                <th>Tarih</th>
                <th>Durum</th>
              </tr>
            </thead>
            <tbody>
              {claims.map((c) => (
                <tr key={c.id}>
                  <td className="text-[#004C3F] font-semibold">#{c.id}</td>
                  <td className="font-mono">{c.claim_number}</td>
                  <td className="max-w-xs truncate">{c.description}</td>
                  <td>{formatDate(c.incident_date)}</td>
                  <td><span className={`badge ${statusClass[c.status] ?? 'badge-gray'}`}>{statusLabel[c.status] ?? c.status}</span></td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
