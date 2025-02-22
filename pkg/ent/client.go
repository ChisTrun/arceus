// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"arceus/pkg/ent/migrate"

	"arceus/pkg/ent/conversation"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Conversation is the client for interacting with the Conversation builders.
	Conversation *ConversationClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Conversation = NewConversationClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:          ctx,
		config:       cfg,
		Conversation: NewConversationClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:          ctx,
		config:       cfg,
		Conversation: NewConversationClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Conversation.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Conversation.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Conversation.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *ConversationMutation:
		return c.Conversation.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// ConversationClient is a client for the Conversation schema.
type ConversationClient struct {
	config
}

// NewConversationClient returns a client for the Conversation from the given config.
func NewConversationClient(c config) *ConversationClient {
	return &ConversationClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `conversation.Hooks(f(g(h())))`.
func (c *ConversationClient) Use(hooks ...Hook) {
	c.hooks.Conversation = append(c.hooks.Conversation, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `conversation.Intercept(f(g(h())))`.
func (c *ConversationClient) Intercept(interceptors ...Interceptor) {
	c.inters.Conversation = append(c.inters.Conversation, interceptors...)
}

// Create returns a builder for creating a Conversation entity.
func (c *ConversationClient) Create() *ConversationCreate {
	mutation := newConversationMutation(c.config, OpCreate)
	return &ConversationCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Conversation entities.
func (c *ConversationClient) CreateBulk(builders ...*ConversationCreate) *ConversationCreateBulk {
	return &ConversationCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *ConversationClient) MapCreateBulk(slice any, setFunc func(*ConversationCreate, int)) *ConversationCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &ConversationCreateBulk{err: fmt.Errorf("calling to ConversationClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*ConversationCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &ConversationCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Conversation.
func (c *ConversationClient) Update() *ConversationUpdate {
	mutation := newConversationMutation(c.config, OpUpdate)
	return &ConversationUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ConversationClient) UpdateOne(co *Conversation) *ConversationUpdateOne {
	mutation := newConversationMutation(c.config, OpUpdateOne, withConversation(co))
	return &ConversationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ConversationClient) UpdateOneID(id uint64) *ConversationUpdateOne {
	mutation := newConversationMutation(c.config, OpUpdateOne, withConversationID(id))
	return &ConversationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Conversation.
func (c *ConversationClient) Delete() *ConversationDelete {
	mutation := newConversationMutation(c.config, OpDelete)
	return &ConversationDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ConversationClient) DeleteOne(co *Conversation) *ConversationDeleteOne {
	return c.DeleteOneID(co.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ConversationClient) DeleteOneID(id uint64) *ConversationDeleteOne {
	builder := c.Delete().Where(conversation.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ConversationDeleteOne{builder}
}

// Query returns a query builder for Conversation.
func (c *ConversationClient) Query() *ConversationQuery {
	return &ConversationQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeConversation},
		inters: c.Interceptors(),
	}
}

// Get returns a Conversation entity by its id.
func (c *ConversationClient) Get(ctx context.Context, id uint64) (*Conversation, error) {
	return c.Query().Where(conversation.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ConversationClient) GetX(ctx context.Context, id uint64) *Conversation {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *ConversationClient) Hooks() []Hook {
	return c.hooks.Conversation
}

// Interceptors returns the client interceptors.
func (c *ConversationClient) Interceptors() []Interceptor {
	return c.inters.Conversation
}

func (c *ConversationClient) mutate(ctx context.Context, m *ConversationMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ConversationCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ConversationUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ConversationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ConversationDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Conversation mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Conversation []ent.Hook
	}
	inters struct {
		Conversation []ent.Interceptor
	}
)
