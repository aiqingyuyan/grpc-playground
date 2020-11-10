package grpc.playground.blog.service

import com.example.protobuf.blog.*
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.*
import java.util.*
import java.util.concurrent.ConcurrentHashMap
import java.util.concurrent.atomic.AtomicInteger

class BlogServiceImpl : BlogServiceGrpcKt.BlogServiceCoroutineImplBase() {

	private var blogCollection: MutableList<Blog>
	private lateinit var authorCount: ConcurrentHashMap<String, Int>
	private lateinit var count: AtomicInteger

	init {
		blogCollection = generateBlogsData()
	}

	private fun generateBlogsData(): MutableList<Blog> {
		val blogList = Collections.synchronizedList(mutableListOf<Blog>())
		count = AtomicInteger(0)
		authorCount = ConcurrentHashMap()

		for (i in 1..10) {
			val blog = Blog
				.newBuilder()
				.setId(count.incrementAndGet())
				.setTitle("Blog $i")
				.setAuthor("Author $i")
				.setText("This is a piece of test text $i")
				.build()

			authorCount[blog.author] = 1
		}

		return blogList
	}

	override suspend fun getBlog(request: GetBlogRequest): GetBlogResponse =
		request.id
			.let { blogId ->
				blogCollection.find { blog ->
					blog.id == blogId
				}
			}
			.let { blog ->
				GetBlogResponse.newBuilder()
					.setBlog(blog)
					.build()
			}

	override suspend fun saveBlogs(requests: Flow<SaveBlogRequest>): SaveBlogResponse {
		requests
			.map { req ->
				Blog.newBuilder()
					.setId(count.incrementAndGet())
					.setTitle(req.blog.title)
					.setAuthor(req.blog.author)
					.setText(req.blog.text)
					.build()
			}
			.collect { blog ->
				blogCollection.add(blog)
			}

		return SaveBlogResponse.newBuilder()
			.setNumberOfBlogs(count.get())
			.build()
	}

	override fun listBlogs(request: ListBlogsRequest): Flow<GetBlogResponse> =
		blogCollection
			.asFlow()
			.onEach { delay(100) }
			.map { blog ->
				GetBlogResponse.newBuilder()
					.setBlog(blog)
					.build()
			}

	override fun getAuthorWithMostBlogsOnSave(requests: Flow<SaveBlogRequest>): Flow<GetAuthorWithMostBlogsResponse> =
		flow {
			requests
				.onEach { delay((0..5).random() * 100L) }
				.map { req ->
					Blog.newBuilder()
						.setId(count.incrementAndGet())
						.setTitle(req.blog.title)
						.setAuthor(req.blog.author)
						.setText(req.blog.text)
						.build()
				}
				.collect { blog ->
					blogCollection.add(blog)

					authorCount.compute(blog.author) { key, _ ->
						if (authorCount[key] == null) {
							1
						} else {
							authorCount[key]!! + 1
						}
					}

					val sorted = authorCount.toSortedMap(compareBy { key ->
						-(authorCount[key] ?: 0)
					})

					emit(GetAuthorWithMostBlogsResponse
						.newBuilder()
						.setAuthor(sorted.firstKey())
						.build()
					)
				}
		}
}
