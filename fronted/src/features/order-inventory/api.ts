import axios, { AxiosError, type AxiosResponse } from 'axios'
import type {
  AddInventoryPayload,
  ApiResponse,
  CreateOrderPayload,
  CreateProductPayload,
  InitInventoryPayload,
  Inventory,
  Order,
  OrderDetail,
  Product,
  StockLog,
} from './types'

const apiBaseURL =
  (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(
    /\/$/,
    ''
  ) || '/api/v1'
const rootBaseURL = apiBaseURL.replace(/\/api\/v1$/, '') || '/'

const api = axios.create({
  baseURL: apiBaseURL,
  timeout: 10000,
})

const rootApi = axios.create({
  baseURL: rootBaseURL,
  timeout: 10000,
})

export class BusinessApiError extends Error {
  code: number

  constructor(message: string, code: number) {
    super(message)
    this.name = 'BusinessApiError'
    this.code = code
  }
}

async function unwrap<T>(
  promise: Promise<AxiosResponse<ApiResponse<T>>>
): Promise<T> {
  const { data } = await promise

  if (data.code !== 0) {
    throw new BusinessApiError(data.message || '请求失败', data.code)
  }

  return data.data as T
}

export function getErrorMessage(error: unknown) {
  if (error instanceof BusinessApiError) return error.message

  if (error instanceof AxiosError) {
    const data = error.response?.data

    if (data && typeof data === 'object') {
      if ('message' in data && typeof data.message === 'string') {
        return data.message
      }
      if ('title' in data && typeof data.title === 'string') {
        return data.title
      }
    }

    if (error.message) return error.message
  }

  if (error instanceof Error && error.message) return error.message

  return '请求失败，请稍后重试'
}

export const queryKeys = {
  health: ['health'] as const,
  products: ['products'] as const,
  product: (id: number) => ['products', id] as const,
  inventory: (productId: number) => ['inventory', productId] as const,
  stockLogsRoot: ['stock-logs'] as const,
  stockLogs: (productId?: number) =>
    ['stock-logs', { productId: productId ?? null }] as const,
  orders: ['orders'] as const,
  order: (id: number) => ['orders', id] as const,
}

export const healthApi = {
  ping: () =>
    unwrap<{ message: string }>(rootApi.get<ApiResponse<{ message: string }>>('/ping')),
}

export const productApi = {
  list: () => unwrap<Product[]>(api.get<ApiResponse<Product[]>>('/products')),
  create: (payload: CreateProductPayload) =>
    unwrap<Product>(api.post<ApiResponse<Product>>('/products', payload)),
  detail: (id: number) =>
    unwrap<Product>(api.get<ApiResponse<Product>>(`/products/${id}`)),
  onSale: (id: number) =>
    unwrap<void>(api.patch<ApiResponse<void>>(`/products/${id}/on-sale`)),
  offSale: (id: number) =>
    unwrap<void>(api.patch<ApiResponse<void>>(`/products/${id}/off-sale`)),
}

export const inventoryApi = {
  init: (payload: InitInventoryPayload) =>
    unwrap<void>(api.post<ApiResponse<void>>('/inventory/init', payload)),
  add: (payload: AddInventoryPayload) =>
    unwrap<void>(api.post<ApiResponse<void>>('/inventory/add', payload)),
  detailByProductId: (productId: number) =>
    unwrap<Inventory>(
      api.get<ApiResponse<Inventory>>(`/inventory/products/${productId}`)
    ),
}

export const stockLogApi = {
  list: (productId?: number) =>
    unwrap<StockLog[]>(
      api.get<ApiResponse<StockLog[]>>('/stock-logs', {
        params: productId ? { product_id: productId } : undefined,
      })
    ),
}

export const orderApi = {
  create: (payload: CreateOrderPayload) =>
    unwrap<Order>(api.post<ApiResponse<Order>>('/orders', payload)),
  list: () => unwrap<Order[]>(api.get<ApiResponse<Order[]>>('/orders')),
  detail: (id: number) =>
    unwrap<OrderDetail>(api.get<ApiResponse<OrderDetail>>(`/orders/${id}`)),
  pay: (id: number) =>
    unwrap<void>(api.patch<ApiResponse<void>>(`/orders/${id}/pay`)),
  finish: (id: number) =>
    unwrap<void>(api.patch<ApiResponse<void>>(`/orders/${id}/finish`)),
  cancel: (id: number) =>
    unwrap<void>(api.patch<ApiResponse<void>>(`/orders/${id}/cancel`)),
}
