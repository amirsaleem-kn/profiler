import { Database } from "../../database";

export interface PromotionalPostMediaDAO { id: number; postID: number; createdBy: string; updatedBy: string; createdAt: number; updatedAt: number; title: string; content: string; isHighlighted: 0 | 1; source: &#39;youtube&#39; | &#39;facebook&#39; | &#39;kisan-network&#39;; mediaType: &#39;image&#39; | &#39;document&#39; | &#39;audio&#39; | &#39;video&#39;; mimeType: string; mediaUrl: string; thumbnailUrl: string; active: 0 | 1;  }

/**
 * @class PromotionalPostMedia
 */
class PromotionalPostMedia {

    public static LIST_QUERY = `SELECT p.id, p.postID, p.createdBy, p.updatedBy, p.createdAt, p.updatedAt, p.title, p.content, p.isHighlighted, p.source, p.mediaType, p.mimeType, p.mediaUrl, p.thumbnailUrl, p.active FROM PromotionalPostMedia p`;

    /**
     * @author Amir Saleem (8 June 2021)
     * @param { Database } db regular db connection
     * @returns { Promise<number> }
     */
    public static insertOne(db: Database, payload: PromotionalPostMediaDAO): Promise<number> {
        return PromotionalPostMedia.insertMany(db, [payload]);
    }

    /**
     * @author Amir Saleem (8 June 2021)
     * @param { Database } db regular db connection
     * @returns { Promise<number> }
     */
    public static async insertMany(db: Database, payload: PromotionalPostMediaDAO[]): Promise<number> {
        const sql = `INSERT INTO PromotionalPostMedia(id, postID, createdBy, updatedBy, createdAt, updatedAt, title, content, isHighlighted, source, mediaType, mimeType, mediaUrl, thumbnailUrl, active) VALUES ?`;
        const values = []
        const result = (await db.executeQuery(sql, values)) as { insertId: number };
        return result.insertId;
    }

}

export default PromotionalPostMedia;
