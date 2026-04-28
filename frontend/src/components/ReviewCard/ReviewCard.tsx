import React from 'react';
import type { Review, User } from '../../data/mockData';
import styles from './ReviewCard.module.css';
import { clsx } from 'clsx';
import { Star } from 'lucide-react';

interface ReviewCardProps {
  review: Review;
  author: User;
}

export const ReviewCard: React.FC<ReviewCardProps> = ({ review, author }) => {
  return (
    <div className={clsx(styles.card, 'glass-panel')}>
      <div className={styles.header}>
        <div className={styles.authorInfo}>
          <img src={author.avatar} alt={author.name} className={styles.avatar} />
          <div>
            <div className={styles.authorName}>{author.name}</div>
            <div className={styles.date}>{review.createdAt}</div>
          </div>
        </div>
        <div className={styles.rating}>
          {[...Array(5)].map((_, i) => (
            <Star
              key={i}
              size={16}
              className={clsx(styles.star, i < review.rating ? styles.starFilled : styles.starEmpty)}
              fill={i < review.rating ? "currentColor" : "none"}
            />
          ))}
        </div>
      </div>
      <div className={styles.text}>{review.text}</div>
    </div>
  );
};
