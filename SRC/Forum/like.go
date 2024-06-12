package Forum

func AddLike(postID int) error {
	_, err := db.Exec("UPDATE Post SET LikeCount = LikeCount + 1 WHERE ID = ?", postID)
	if err != nil {
		return err
	}
	return nil
}

func RemoveLike(postID int) error {
	_, err := db.Exec("UPDATE Post SET LikeCount = LikeCount - 1 WHERE ID = ?", postID)
	if err != nil {
		return err
	}
	return nil
}
