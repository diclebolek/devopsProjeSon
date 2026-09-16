'use client';
import { useEffect, useState } from 'react';
import { apiService } from '@/lib/api';
import { formatDate } from '@/lib/dates';
import { Quote } from '@/types';

const statusClass: Record<string, string> = {
  draft: 'badge-gray',
  sent: 'badge-yellow',
  accepted: 'badge-green',
  rejected: 'badge-red',
};

const statusLabel: Record<string, string> = {
  draft: 'Taslak',
  sent: 'Gönderildi',
  accepted: 'Kabul',
  rejected: 'Red',
};

export default function QuotesPage() {
  const [quotes, setQuotes] = useState<Quote[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    apiService.getMyQuotes()
      .then((res) => setQuotes(res.data as Quote[]))
      .catch(() => setError('Teklifler yüklenirken bir hata oluştu.'))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-[#004C3F]">Tekliflerim</h1>
          <p className="text-[#677471] text-sm mt-1">Talep ettiğiniz sigorta teklifleri</p>
        </div>
        <span className="badge badge-green">{quotes.length} Teklif</span>
      </div>

      <div className="card overflow-hidden">
        {loading ? (
          <div className="flex justify-center py-16"><div className="spinner" /></div>
        ) : error ? (
          <div className="text-center py-16 text-red-500">{error}</div>
        ) : quotes.length === 0 ? (
          <div className="text-center py-16 text-[#677471]">
            <span className="text-5xl block mb-4">💼</span>
            <p>Henüz teklif talebiniz bulunmuyor.</p>
          </div>
        ) : (
          <table className="dash-table">
            <thead>
              <tr>
                <th>No</th>
                <th>Hizmet</th>
                <th>Talep Tarihi</th>
                <th>Durum</th>
              </tr>
            </thead>
            <tbody>
              {quotes.map((q) => (
                <tr key={q.id}>
                  <td className="text-[#004C3F] font-semibold">{q.reference_code || `#${q.id}`}</td>
                  <td>{q.product_type}</td>
                  <td>{formatDate(q.created_at)}</td>
                  <td><span className={`badge ${statusClass[q.status] ?? 'badge-gray'}`}>{statusLabel[q.status] ?? q.status}</span></td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
