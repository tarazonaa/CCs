
export interface Drawing {
  id: string;
  imageData: string;
  prediction?: string;
  confidence?: number;
  timestamp: Date;
}
