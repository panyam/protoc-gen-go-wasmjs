import { BookStats as BookStatsInterface, UserActivity as UserActivityInterface, GetBookStatsRequest as GetBookStatsRequestInterface, GetBookStatsResponse as GetBookStatsResponseInterface } from "./analytics_interfaces";


/**
 * BookStats provides analytics for books
 */
export class BookStats implements BookStatsInterface {
  base?: BaseMessage;
  bookId: string = "";
  viewCount: number = 0;
  checkoutCount: number = 0;
  averageRating: number = 0;
  reviewCount: number = 0;
}


/**
 * UserActivity tracks user behavior
 */
export class UserActivity implements UserActivityInterface {
  base?: BaseMessage;
  userId: string = "";
  activityType: string = "";
  bookId: string = "";
  description: string = "";
}


/**
 * GetBookStatsRequest for analytics
 */
export class GetBookStatsRequest implements GetBookStatsRequestInterface {
  metadata?: Metadata;
  bookId: string = "";
  dateRange: string = "";
}


/**
 * GetBookStatsResponse for analytics
 */
export class GetBookStatsResponse implements GetBookStatsResponseInterface {
  metadata?: Metadata;
  stats?: BookStats;
  error?: ErrorInfo;
}

