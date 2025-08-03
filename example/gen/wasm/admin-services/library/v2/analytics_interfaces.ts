// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * BookStats provides analytics for books
 */
export interface BookStats {
  base?: BaseMessage;
  bookId: string;
  viewCount: number;
  checkoutCount: number;
  averageRating: number;
  reviewCount: number;
}


/**
 * UserActivity tracks user behavior
 */
export interface UserActivity {
  base?: BaseMessage;
  userId: string;
  activityType: string;
  bookId: string;
  description: string;
}


/**
 * GetBookStatsRequest for analytics
 */
export interface GetBookStatsRequest {
  metadata?: Metadata;
  bookId: string;
  dateRange: string;
}


/**
 * GetBookStatsResponse for analytics
 */
export interface GetBookStatsResponse {
  metadata?: Metadata;
  stats?: BookStats;
  error?: ErrorInfo;
}

