import { useState, useEffect, useCallback } from 'react';
import { apiRequest } from '../api/client';

export const usePagination = (initialPage = 1, initialPerPage = 10) => {
  const [currentPage, setCurrentPage] = useState(initialPage);
  const [perPage, setPerPage] = useState(initialPerPage);
  const [results, setResults] = useState([]);
  const [pagination, setPagination] = useState({
    current_page: 1,
    per_page: 10,
    total_items: 0,
    total_pages: 0,
    has_next: false,
    has_prev: false
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchResults = useCallback(async (page = currentPage, itemsPerPage = perPage) => {
    setLoading(true);
    setError(null);
    
    try {
      const { ok, data } = await apiRequest(
        `/results?page=${page}&per_page=${itemsPerPage}`
      );
      
      if (ok && data.data) {
        setResults(data.data.data || []);
        const paginationData = data.data.pagination || {};
        
        setPagination({
          current_page: paginationData.current_page || page,
          per_page: itemsPerPage,
          total_items: paginationData.total_items || 0,
          total_pages: Math.max(1, paginationData.total_pages || 1),
          has_next: paginationData.has_next || false,
          has_prev: paginationData.has_prev || false
        });
      } else {
        setError(data.error || 'Error al cargar resultados');
        setResults([]);
        setPagination({
          current_page: 1,
          per_page: itemsPerPage,
          total_items: 0,
          total_pages: 1,
          has_next: false,
          has_prev: false
        });
      }
    } catch (err) {
      setError('Error de conexión');
      setResults([]);
      setPagination({
        current_page: 1,
        per_page: itemsPerPage,
        total_items: 0,
        total_pages: 1,
        has_next: false,
        has_prev: false
      });
    } finally {
      setLoading(false);
    }
  }, [currentPage, perPage]);

  const handlePageChange = useCallback((newPage) => {
    if (newPage >= 1 && newPage <= Math.max(1, pagination.total_pages)) {
      setCurrentPage(newPage);
    }
  }, [pagination.total_pages]);

  const handlePerPageChange = useCallback((newPerPage) => {
    setPerPage(newPerPage);
    setCurrentPage(1);
  }, []);

  const handleDelete = useCallback(async (id) => {
    const { ok } = await apiRequest(`/results/${id}`, { method: "DELETE" });
    
    if (ok) {
      const newTotal = Math.max(0, pagination.total_items - 1);
      const newTotalPages = Math.max(1, Math.ceil(newTotal / perPage));
      
      if (currentPage > newTotalPages && newTotalPages > 0) {
        setCurrentPage(newTotalPages);
      } else {
        fetchResults(currentPage, perPage);
      }
    }
    return ok;
  }, [pagination.total_items, perPage, currentPage, fetchResults]);

  const refresh = useCallback(() => {
    fetchResults(currentPage, perPage);
  }, [fetchResults, currentPage, perPage]);

  useEffect(() => {
    fetchResults(currentPage, perPage);
  }, [currentPage, perPage, fetchResults]);

  const totalItems = pagination.total_items;
  const totalPages = Math.max(1, pagination.total_pages);
    const startItem = totalItems === 0 ? 0 : Math.min((currentPage - 1) * perPage + 1, totalItems);
  const endItem = totalItems === 0 ? 0 : Math.min(currentPage * perPage, totalItems);

  return {
    results,
    pagination,
    loading,
    error,
    currentPage,
    perPage,
    handlePageChange,
    handlePerPageChange,
    handleDelete,
    refresh,
    hasNext: pagination.has_next && currentPage < totalPages,
    hasPrev: pagination.has_prev && currentPage > 1,
    totalItems,
    totalPages,
    startItem,
    endItem
  };
};