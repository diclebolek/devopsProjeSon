'use client';
import { useEffect, useState } from 'react';
import { apiService } from '@/lib/api';
import { formatDate } from '@/lib/dates';
import { Payment } from '@/types';

const statusClass: Record<string, string> = {
  paid: 'badge-green',
  pending: 'badge-yellow',
  overdue: 'badge-red',
};

const statusLabel: Record<string, string> = {
  paid: 'Ödendi',
  pending: 'Bekliyor',
  overdue: 'Gecikmiş',
};

export default function PaymentsPage() {
  const [payments, setPayments] = useState<Payment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    apiService.getMyPayments()
      .then((res) => setPayments(res.data as Payment[]))
      .catch(() => setError('Ödemeler yüklenirken bir hata oluştu.'))
      .finally(() => setLoading(false));
  }, []);

  const total = payments
    .filter((p) => p.status === 'paid')
    .reduce((sum, p) => sum + parseFloat(p.amount), 0);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-[#004C3F]">Ödemelerim</h1>
          <p className="text-[#677471] text-sm mt-1">Geçmiş ve bekleyen prim ödemeleriniz</p>
        </div>
        <div className="text-right">
          <p className="text-xs text-[#677471]">Toplam Ödenen</p>
          <p className="text-xl font-bold text-[#028835]">{total.toFixed(2)} ₺</p>
        </div>
      </div>

      <div className="card overflow-hidden">
        {loading ? (
          <div className="flex justify-center py-16"><div className="spinner" /></div>
        ) : error ? (
          <div className="text-center py-16 text-red-500">{error}</div>
        ) : payments.length === 0 ? (
          <div className="text-center py-16 text-[#677471]">
            <span className="text-5xl block mb-4">💳</span>
            <p>Ödeme geçmişiniz bulunmuyor.</p>
          </div>
        ) : (
          <table className="dash-table">
            <thead>
              <tr>
                <th>No</th>
                <th>Açıklama</th>
                <th>Tutar</th>
                <th>Vade Tarihi</th>
                <th>Durum</th>
              </tr>
            </thead>
            <tbody>
              {payments.map((p) => (
                <tr key={p.id}>
                  <td className="text-[#004C3F] font-semibold">#{p.id}</td>
                  <td>{p.description}</td>
                  <td className="font-semibold text-[#004C3F]">{parseFloat(p.amount).toFixed(2)} ₺</td>
                  <td>{formatDate(p.due_date)}</td>
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
